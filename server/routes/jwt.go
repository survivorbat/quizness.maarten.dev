package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/survivorbat/qq.maarten.dev/server/services"
	"net/http"
)

const bearerSchema = "Bearer"

type JwtHandler struct {
	JwtService services.IJwtService
}

func (j *JwtHandler) AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) <= len(bearerSchema) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenString := authHeader[len(bearerSchema):]
		token, err := j.JwtService.ValidateToken(tokenString)

		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("user", token.Claims.(services.QQClaims).UserID)
	}
}
