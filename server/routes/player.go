package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/survivorbat/qq.maarten.dev/server/domain"
	"github.com/survivorbat/qq.maarten.dev/server/services"
	"net/http"
)

type PlayerHandler struct {
	PlayerService services.PlayerService
	GameService   services.GameService
}

// Get godoc
//
//	@Summary	Fetch this game's players
//	@Tags		Player
//	@Accept		json
//	@Produce	json
//	@Param		id	path	string			true	"ID of the game"
//	@Success	200	{array}	[]domain.Player	"This game's players"
//	@Failure	400	"Invalid uuid"
//	@Failure	403	"You can only view your own game's players"
//	@Failure	500	"Internal Server Error"
//	@Router		/api/v1/games/{id}/players [get]
//	@Security	JWT
func (g *PlayerHandler) Get(c *gin.Context) {
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

	// Prevent users from viewing other people's games
	if game.Quiz.CreatorID.String() != authID {
		logrus.Errorf("Creator is %s not %s", game.Quiz.CreatorID, authID)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	games, err := g.PlayerService.GetByGame(gameID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get by game")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, games)
}

// Post godoc
//
//	@Summary	Create a new player for this game
//	@Tags		Player
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string			true	"ID of the game"
//	@Success	200	{object}	domain.Player	"The new player"
//	@Failure	400	"Invalid uuid"
//	@Failure	404	"Game not found"
//	@Failure	409	"Game full"
//	@Failure	500	"Internal Server Error"
//	@Router		/api/v1/games/{id}/players [post]
func (g *PlayerHandler) Post(c *gin.Context) {
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

	// Don't show them it exists
	if !game.IsInProgress() {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if uint(len(game.Players)) >= game.PlayerLimit {
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	player := &domain.Player{Game: game}
	if err := g.PlayerService.Create(player); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, player)
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
	id := c.Param("id")

	playerID, err := uuid.Parse(id)
	if err != nil {
		logrus.WithError(err).Error("UUID error")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	player, err := g.PlayerService.GetByID(playerID)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if err := g.PlayerService.Delete(player); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, player)
}
