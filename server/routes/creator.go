package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/survivorbat/qq.maarten.dev/server/services"
	"net/http"
)

type CreatorHandler struct {
	CreatorService services.ICreatorService
}

func (g *CreatorHandler) GetWithID(c *gin.Context) {
	id := c.Param("id")
	authID := c.GetString("user")

	// Users may only fetch their own data
	if id != authID {
		logrus.Errorf("User is not the authenticated one")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	creator, err := g.CreatorService.GetByID(uuid.MustParse(authID))
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch creator by ID")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, creator)
}
