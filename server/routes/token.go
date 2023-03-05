package routes

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/survivorbat/qq.maarten.dev/server/services"
	"github.com/zalando/gin-oauth2/google"
	"golang.org/x/oauth2"
	"io"
	"net/http"
)

type TokenHandler struct {
	CreatorService services.ICreatorService
	JwtService     services.IJwtService
	AuthConfig     *oauth2.Config
}

type tokenInput struct {
	Code string `json:"code"`
}

func (a *TokenHandler) CreateToken(c *gin.Context) {
	var input *tokenInput
	if err := c.ShouldBindJSON(&input); err != nil {
		logrus.WithError(err).Error("Failed to bind json")
		c.AbortWithStatus(http.StatusBadRequest)
	}

	tok, err := a.AuthConfig.Exchange(oauth2.NoContext, input.Code)
	if err != nil {
		logrus.WithError(err).Error("Failed to exchange token")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	client := a.AuthConfig.Client(context.TODO(), tok)
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
	user, err := a.CreatorService.GetOrCreate(googleUser.Sub)
	if err != nil {
		logrus.WithError(err).Error("Failed to create or get user")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	token, err := a.JwtService.GenerateToken(user.ID)
	if err != nil {
		logrus.WithError(err).Error("Failed to generate token")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Header("token", token)
	c.Status(http.StatusOK)
}
