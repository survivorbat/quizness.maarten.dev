package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/survivorbat/qq.maarten.dev/server/services"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type GameConnectionHandler struct {
	GameService services.GameService
}

// Get godoc
//
//	@Summary	Connect to this game using a websocket
//	@Tags		Game
//	@Accept		json
//	@Produce	json
//	@Param		id		path	string	true	"ID of the game"
//	@Param		player	path	string	true	"ID of the player"
//	@Success	200		"An established connection"
//	@Failure	400		"Invalid uuid"
//	@Failure	403		"Player is not in game"
//	@Failure	404		"Game not found"
//	@Failure	404		"Game is not open for joining"
//	@Failure	500		"Internal Server Error"
//	@Router		/api/v1/games/{id}/players/{player}/connection [get]
func (g *GameConnectionHandler) Get(c *gin.Context) {
	gameParam := c.Param("id")
	gameID, err := uuid.Parse(gameParam)
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

	if !game.Players.Contains(playerID) {
		logrus.Error("Player is not in game")
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logrus.WithError(err).Error("Game can not be joined")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	defer ws.Close()

	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			logrus.WithError(err).Error("Failed to read")
			break
		}

		logrus.Info(mt, message)
	}
}

// GetCreator godoc
//
//	@Summary	Connect to this game using a websocket as a creator
//	@Tags		Game
//	@Accept		json
//	@Produce	json
//	@Param		id	path	string	true	"ID of the game"
//	@Success	200	"An established connection"
//	@Failure	400	"Invalid uuid"
//	@Failure	403	"Not your game"
//	@Failure	404	"Game not found"
//	@Failure	404	"Game is not started"
//	@Failure	500	"Internal Server Error"
//	@Router		/api/v1/games/{id}/connection [get]
//	@Security	JWT
func (g *GameConnectionHandler) GetCreator(c *gin.Context) {
	authID := c.GetString("user")
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

	if game.Quiz.CreatorID != uuid.MustParse(authID) {
		logrus.Error("Not your game")
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logrus.WithError(err).Error("Game can not be joined")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	defer ws.Close()

	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			logrus.WithError(err).Error("Failed to read")
			break
		}

		logrus.Info(mt, message)
	}
}