package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/survivorbat/qq.maarten.dev/server/domain"
	"github.com/survivorbat/qq.maarten.dev/server/routes/inputs"
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
//	@Tags		Quiz
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
		logrus.WithError(err).Error("Failed to get by quiz")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, games)
}

// Post godoc
//
//	@Summary	Create a new game for this quiz
//	@Tags		Quiz
//	@Accept		json
//	@Produce	json
//	@Param		input	body		inputs.Game	true	"Your game"
//	@Success	200		{object}	domain.Game	"The new game"
//	@Failure	400		{object}	any			"Invalid uuid"
//	@Failure	400		{object}	any			"You already have a game started"
//	@Failure	403		{object}	any			"You can only create games on your own quiz"
//	@Failure	500		{object}	any			"Internal Server Error"
//	@Router		/api/v1/quizzes/:id/games [post]
//	@Security	JWT
func (g *GameHandler) Post(c *gin.Context) {
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

	// TODO: Add only 1 check

	var input *inputs.Game
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	game := &domain.Game{
		QuizID:      quizID,
		PlayerLimit: input.PlayerLimit,
	}

	if err := g.GameService.Create(game); err != nil {
		logrus.WithError(err).Error("Failed to create")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, game)
}
