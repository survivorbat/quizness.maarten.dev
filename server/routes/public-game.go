package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/survivorbat/qq.maarten.dev/server/routes/inputs"
	"github.com/survivorbat/qq.maarten.dev/server/routes/outputs"
	"github.com/survivorbat/qq.maarten.dev/server/services"
	"net/http"
)

type PublicGameHandler struct {
	GameService services.GameService
}

// Get godoc
//
//	@Summary	Get the current question
//	@Tags		Game
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string									true	"ID of the game"
//	@Success	200	{object}	outputs.OutputMultipleChoiceQuestion	"The current question"
//	@Failure	400	"Invalid uuid"
//	@Failure	404	"Game not found"
//	@Failure	404	"Game not active"
//	@Failure	500	"Internal Server Error"
//	@Router		/api/v1/games/{id}/questions/current [get]
func (g *PublicGameHandler) Get(c *gin.Context) {
	gameParam := c.Param("id")
	gameID, err := uuid.Parse(gameParam)
	if err != nil {
		logrus.WithError(err).Error("UUID error")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	game, err := g.GameService.GetByID(gameID)
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch game")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	question, ok := game.GetCurrentQuestion()
	if !ok {
		logrus.WithError(err).Error("Failed to fetch question")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	output, err := outputs.NewPublicQuestion(question)
	if err != nil {
		logrus.WithError(err).Error("Failed to create output")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.Data(http.StatusOK, "application/json", output)
}

// Patch godoc
//
//	@Summary	Answer a question
//	@Tags		Game
//	@Accept		json
//	@Produce	json
//	@Param		id			path		string				true	"ID of the game"
//	@Param		question	path		string				true	"ID of the question"
//	@Param		player		path		string				true	"ID of the player"
//	@Param		input		body		inputs.Answer		true	"Your answer"
//	@Success	200			{object}	domain.GameAnswer	"The new answer"
//	@Failure	400			"Invalid uuid"
//	@Failure	403			"You can only answer a question while it's active"
//	@Failure	404			"Quiz not found"
//	@Failure	404			"Question not found"
//	@Failure	500			"Internal Server Error"
//	@Router		/api/v1/games/{id}/questions/{question}/answers/{player} [patch]
func (g *PublicGameHandler) Patch(c *gin.Context) {
	gameParam := c.Param("id")
	gameID, err := uuid.Parse(gameParam)
	if err != nil {
		logrus.WithError(err).Error("UUID error")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	questionParam := c.Param("question")
	questionID, err := uuid.Parse(questionParam)
	if err != nil {
		logrus.WithError(err).Error("UUID error")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	playerParam := c.Param("player")
	playerID, err := uuid.Parse(playerParam)
	if err != nil {
		logrus.WithError(err).Error("UUID error")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	game, err := g.GameService.GetByID(gameID)
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch game")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	var input *inputs.Answer
	if err := c.ShouldBindJSON(&input); err != nil {
		logrus.WithError(err).Error("Validation error")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := g.GameService.AnswerQuestion(game, questionID, playerID, input.OptionID); err != nil {
		logrus.WithError(err).Error("Answer error")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
