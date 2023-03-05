package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/survivorbat/qq.maarten.dev/server/services"
	"net/http"
)

type QuizHandler struct {
	QuizService services.IQuizService
}

func (g *QuizHandler) Get(c *gin.Context) {
	authID := c.GetString("user")

	quizzes, err := g.QuizService.GetByCreator(uuid.MustParse(authID))
	if err != nil {
		logrus.WithError(err).Error("Failed to get by creator")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, quizzes)
}
