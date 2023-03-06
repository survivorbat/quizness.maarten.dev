package routes

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/survivorbat/qq.maarten.dev/server/services"
	"net/http"
)

const bearerSchema = "Bearer"

type JwtHandler struct {
	JwtService services.IJwtService
}

func (j *JwtHandler) JwtGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) <= len(bearerSchema) {
			logrus.Error("Authorization header is wrong")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenString := authHeader[len(bearerSchema)+1:]
		token, err := j.JwtService.ValidateToken(tokenString)

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
//	@Failure	200		{object}	any					"Token in the header"
//	@Failure	500	{object}	any				"Internal Server Error"
//	@Router		/api/v1/tokens [put]
//	@Security	JWT
func (j *JwtHandler) Refresh(c *gin.Context) {
	authID := c.GetString("user")

	token, err := j.JwtService.GenerateToken(authID)
	if err != nil {
		logrus.WithError(err).Error("Failed to generate token")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Header("token", token)
	c.Status(http.StatusOK)
}
