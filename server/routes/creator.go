package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/survivorbat/qq.maarten.dev/server/services"
	"net/http"
)

type CreatorHandler struct {
	CreatorService services.CreatorService
}

// GetWithID godoc
//
//	@Summary	Fetch your account's data
//	@Tags		Creator
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	domain.Creator	"The creator"
//	@Failure	500	{object}	any				"Internal Server Error"
//	@Router		/api/v1/creators/self [get]
//	@Security	JWT
func (g *CreatorHandler) GetWithID(c *gin.Context) {
	authID := c.GetString("user")

	creator, err := g.CreatorService.GetByID(uuid.MustParse(authID))
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch creator by ID")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, creator)
}
