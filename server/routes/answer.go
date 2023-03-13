package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type AnswerHandler struct {
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
func (g *AnswerHandler) Patch(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
	return
}
