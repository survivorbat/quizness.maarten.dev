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

type GameControlHandler struct {
	QuizService services.QuizService
	GameService services.GameService
}

// Get godoc
//
//	@Summary	Fetch this quiz' games
//	@Tags		Quiz
//	@Accept		json
//	@Produce	json
//	@Param		id	path	string			true	"ID of the quiz"
//	@Success	200	{array}	[]domain.Game	"This quiz' games"
//	@Failure	400	"Invalid uuid"
//	@Failure	403	"You can only view your own games"
//	@Failure	500	"Internal Server Error"
//	@Router		/api/v1/quizzes/{id}/games [get]
//	@Security	JWT
func (g *GameControlHandler) Get(c *gin.Context) {
	authID := c.GetString("user")
	id := c.Param("id")

	quizID, err := uuid.Parse(id)
	if err != nil {
		logrus.WithError(err).Error("UUID error")
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
//	@Param		id		path		string		true	"ID of the quiz"
//	@Param		input	body		inputs.Game	true	"Your game"
//	@Success	200		{object}	domain.Game	"The new game"
//	@Failure	400		"Invalid uuid"
//	@Failure	403		"You can only create games on your own quiz"
//	@Failure	500		"Internal Server Error"
//	@Router		/api/v1/quizzes/{id}/games [post]
//	@Security	JWT
func (g *GameControlHandler) Post(c *gin.Context) {
	authID := c.GetString("user")
	id := c.Param("id")

	quizID, err := uuid.Parse(id)
	if err != nil {
		logrus.WithError(err).Error("UUID error")
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

	var input *inputs.Game
	if err := c.ShouldBindJSON(&input); err != nil {
		logrus.WithError(err).Error("Validation error")
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

// Patch godoc
//
//	@Summary	Perform actions on a game
//	@Tags		Game
//	@Accept		json
//	@Produce	json
//	@Param		id		path		string		true	"ID of the game"
//	@Param		action	query		string		true	"Action to perform"	Enums(start, finish, next)
//	@Success	200		{object}	domain.Game	"The updated game"
//	@Failure	400		"Invalid uuid"
//	@Failure	400		"Game is not in a valid state"
//	@Failure	400		"Unknown action"
//	@Failure	403		"You can only change games in your own quiz"
//	@Failure	404		"Not found"
//	@Failure	409		"You already have a game started"
//	@Failure	500		"Internal Server Error"
//	@Router		/api/v1/games/{id} [patch]
//	@Security	JWT
func (g *GameControlHandler) Patch(c *gin.Context) {
	authID := c.GetString("user")
	id := c.Param("id")
	action := c.Query("action")

	gameID, err := uuid.Parse(id)
	if err != nil {
		logrus.WithError(err).Error("UUID error")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	game, err := g.GameService.GetByID(gameID)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Prevent users from viewing other people's games
	if game.Quiz.CreatorID.String() != authID {
		logrus.Errorf("Creator is %s not %s", game.Quiz.CreatorID, authID)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	switch action {
	case "start":
		if game.Quiz.HasGameInProgress() {
			c.AbortWithStatus(http.StatusConflict)
			return
		}

		if err := g.GameService.Start(game); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

	case "finish":
		if err := g.GameService.Finish(game); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

	case "next":
		// TODO: Differentiate errors
		if err := g.GameService.Next(game); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	default:
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, game)
}

// Delete godoc
//
//	@Summary	Delete a game
//	@Tags		Game
//	@Accept		json
//	@Produce	json
//	@Param		id	path	string	true	"ID of the game"
//	@Success	200	"The deleted game"
//	@Failure	400	"Invalid uuid"
//	@Failure	403	"You can only delete games in your own quiz"
//	@Failure	404	"Not found"
//	@Failure	409	"This game has started"
//	@Failure	500	"Internal Server Error"
//	@Router		/api/v1/games/{id} [delete]
//	@Security	JWT
func (g *GameControlHandler) Delete(c *gin.Context) {
	authID := c.GetString("user")
	id := c.Param("id")

	gameID, err := uuid.Parse(id)
	if err != nil {
		logrus.WithError(err).Error("UUID error")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	game, err := g.GameService.GetByID(gameID)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Prevent users from deleting other people's games
	if game.Quiz.CreatorID.String() != authID {
		logrus.Errorf("Creator is %s not %s", game.Quiz.CreatorID, authID)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	if err := g.GameService.Delete(game); err != nil {
		logrus.WithError(err).Error("Failed to delete")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, game)
}
