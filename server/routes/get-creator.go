package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/survivorbat/qq.maarten.dev/server/services"
	"net/http"
)

type GetCreatorHandler struct {
	CreatorService services.ICreatorService
}

func (g *GetCreatorHandler) GetWithID(c *gin.Context) {
	id := c.Param("id")
	authID := c.GetString("user")

	// Users may only fetch their own data
	if id != authID {
		c.Status(http.StatusUnauthorized)
		return
	}

	creator, err := g.CreatorService.GetByID(uuid.MustParse(authID))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, creator)
}
