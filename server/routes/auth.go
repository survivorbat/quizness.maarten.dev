package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/survivorbat/qq.maarten.dev/server/services"
	"github.com/zalando/gin-oauth2/google"
	"net/http"
)

type AuthHandler struct {
	CreatorService services.ICreatorService
	JwtService     services.IJwtService
}

func (a *AuthHandler) Handle(c *gin.Context) {
	val := c.MustGet("user")
	res, _ := val.(google.User)

	// Ensure user exists
	user, err := a.CreatorService.GetOrCreate(res.Sub)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	token, err := a.JwtService.GenerateToken(user.ID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Header("token", token)
	c.Redirect(http.StatusFound, "/creator")
}
