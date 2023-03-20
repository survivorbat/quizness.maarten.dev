package routes

import (
	"context"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/survivorbat/qq.maarten.dev/server/inputs"
	"github.com/survivorbat/qq.maarten.dev/server/services"
	"github.com/zalando/gin-oauth2/google"
	"golang.org/x/oauth2"
	"io"
	"net/http"
)

const bearerSchema = "Bearer"

type TokenHandler struct {
	CreatorService services.CreatorService
	JwtService     services.JwtService
	AuthConfig     *oauth2.Config
}

// CreateToken godoc
//
//	@Summary	Create a new authentication token using OAuth
//	@Tags		Token
//	@Accept		json
//	@Produce	json
//	@Param		code	body	inputs.Token	true	"Your OAuth code"
//	@Failure	200		"Token in the header"
//	@Failure	400		"Malformed input"
//	@Failure	401		"Failed to authenticate you"
//	@Failure	500		"Internal Server Error"
//	@Router		/api/v1/tokens [post]
func (t *TokenHandler) CreateToken(c *gin.Context) {
	var input *inputs.Token
	if err := c.ShouldBindJSON(&input); err != nil {
		logrus.WithError(err).Error("Validation error")
		c.AbortWithStatus(http.StatusBadRequest)
	}

	tok, err := t.AuthConfig.Exchange(oauth2.NoContext, input.Code)
	if err != nil {
		logrus.WithError(err).Error("Failed to exchange token")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	client := t.AuthConfig.Client(context.TODO(), tok)
	email, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch user info")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	defer email.Body.Close()
	data, err := io.ReadAll(email.Body)
	if err != nil {
		logrus.WithError(err).Error("Failed to read body")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var googleUser google.User
	if err := json.Unmarshal(data, &googleUser); err != nil {
		logrus.WithError(err).Error("Failed to read user data")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Ensure user exists
	user, err := t.CreatorService.GetOrCreate(googleUser.Sub)
	if err != nil {
		logrus.WithError(err).Error("Failed to create or get user")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	token, err := t.JwtService.GenerateToken(user.ID.String())
	if err != nil {
		logrus.WithError(err).Error("Failed to generate token")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Header("token", token)
	c.Status(http.StatusOK)
}

func (t *TokenHandler) JwtGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			authHeader = c.GetHeader("Sec-Websocket-Protocol")
		}

		if len(authHeader) <= len(bearerSchema) {
			logrus.Error("Authorization header is wrong")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenString := authHeader[len(bearerSchema)+1:]
		token, err := t.JwtService.ValidateToken(tokenString)

		if err != nil {
			logrus.WithError(err).Error("Failed to validate token")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			logrus.Error("Token is invalid")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("user", token.Claims.(jwt.MapClaims)["userID"])
	}
}

// Refresh godoc
//
//	@Summary	Refresh your authentication token
//	@Tags		Token
//	@Accept		json
//	@Produce	json
//	@Failure	200	"Token in the header"
//	@Failure	500	"Internal Server Error"
//	@Router		/api/v1/tokens [put]
//	@Security	JWT
func (t *TokenHandler) Refresh(c *gin.Context) {
	authID := c.GetString("user")

	token, err := t.JwtService.GenerateToken(authID)
	if err != nil {
		logrus.WithError(err).Error("Failed to generate token")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Header("token", token)
	c.Status(http.StatusOK)
}
