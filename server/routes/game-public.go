package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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

// GetByCode godoc
//
//	@Summary	Get a game by its code
//	@Tags		Game
//	@Accept		json
//	@Produce	json
//	@Param		code	query		string				true	"Code of the game"
//	@Success	200		{object}	outputs.OutputGame	"The game ID"
//	@Failure	403		"Can only be used for filtering on codes"
//	@Failure	404		"Game not found"
//	@Failure	500		"Internal Server Error"
//	@Router		/api/v1/games [get]
func (g *PublicGameHandler) GetByCode(c *gin.Context) {
	code := c.Query("code")

	if code == "" {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	game, err := g.GameService.GetByCode(code)
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch game")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, outputs.NewPublicGame(game))
}
