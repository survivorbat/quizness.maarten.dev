package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type PlayerHandler struct {
}

// Get godoc
//
//	@Summary	Fetch this quiz' players
//	@Tags		Player
//	@Accept		json
//	@Produce	json
//	@Param		id	path	string			true	"ID of the quiz"
//	@Success	200	{array}	[]domain.Player	"This quiz' players"
//	@Failure	400	"Invalid uuid"
//	@Failure	403	"You can only view your own quiz' players"
//	@Failure	500	"Internal Server Error"
//	@Router		/api/v1/quizzes/{id}/players [get]
//	@Security	JWT
func (g *PlayerHandler) Get(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
	return
}

// Post godoc
//
//	@Summary	Create a new player for this quiz
//	@Tags		Player
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string			true	"ID of the quiz"
//	@Success	200	{object}	domain.Player	"The new player"
//	@Failure	400	"Invalid uuid"
//	@Failure	403	"You can only join a quiz that is active"
//	@Failure	404	"Player not found"
//	@Failure	500	"Internal Server Error"
//	@Router		/api/v1/quizzes/{id}/players [post]
func (g *PlayerHandler) Post(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
	return
}

// Delete godoc
//
//	@Summary	Delete a player
//	@Tags		Player
//	@Accept		json
//	@Produce	json
//	@Param		id	path	string	true	"ID of the player"
//	@Success	200	"The deleted player"
//	@Failure	400	"Invalid uuid"
//	@Failure	404	"Not found"
//	@Failure	500	"Internal Server Error"
//	@Router		/api/v1/players/{id} [delete]
func (g *PlayerHandler) Delete(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
	return
}
