package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/survivorbat/qq.maarten.dev/server/services"
	"net/http"
)

type GameHandler struct {
	QuizService services.QuizService
	GameService services.GameService
}

// Get godoc
//
//	@Summary	Fetch this quiz' games
//	@Tags		Game
//	@Accept		json
//	@Produce	json
//	@Success	200	{array}		[]domain.Game	"This quiz' games"
//	@Failure	400	{object}	any				"Invalid uuid"
//	@Failure	403	{object}	any				"You can only view your own games"
//	@Failure	500	{object}	any				"Internal Server Error"
//	@Router		/api/v1/quizzes/:id/games [get]
//	@Security	JWT
func (g *GameHandler) Get(c *gin.Context) {
	authID := c.GetString("user")
	id := c.Param("id")

	quizID, err := uuid.Parse(id)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	quiz, err := g.QuizService.GetByID(quizID)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Prevent users from viewing other people's games
	if quiz.CreatorID.String() != authID {
		logrus.Errorf("Creator is %s not %s", quiz.CreatorID, authID)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	games, err := g.GameService.GetByQuiz(quizID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get by creator")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, games)
}
