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

// Get godoc
//
//	@Summary	Fetch your quizzes
//	@Tags		Quiz
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	[]domain.Quiz	"Your quizzes"
//	@Failure	500	{object}	any				"Internal Server Error"
//	@Router		/api/v1/quizzes [get]
//	@Security	JWT
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
