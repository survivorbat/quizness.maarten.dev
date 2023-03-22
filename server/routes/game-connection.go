package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/survivorbat/qq.maarten.dev/server/coordinator"
	"github.com/survivorbat/qq.maarten.dev/server/services"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type GameConnectionHandler struct {
	GameService    services.GameService
	PlayerService  services.PlayerService
	CreatorService services.CreatorService
	Coordinator    coordinator.GameCoordinator
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
//	@Failure	400		"Invalid websocket headers"
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

	player, err := g.PlayerService.GetByID(playerID)
	if err != nil {
		logrus.WithError(err).Error("How even")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	logrus.Infof("Opening websocket for player %s in game %s", playerID, gameID)
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logrus.WithError(err).Error("Game can not be joined")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	defer ws.Close()

	finish := make(chan struct{})

	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("Had to recover from a panic: %v", err)
		}

		g.Coordinator.UnsubscribePlayer(gameID, player)
	}()

	g.Coordinator.SubscribePlayer(gameID, player, func(message *coordinator.BroadcastMessage) {
		if err := ws.WriteJSON(message); err != nil {
			logrus.WithError(err).Error("Failed to write JSON")
		}

		if message.Type == coordinator.FinishGameType {
			finish <- struct{}{}
		}
	})

	for {
		select {
		default:
			var result *coordinator.PlayerMessage
			if err := ws.ReadJSON(&result); err != nil {
				continue
			}

			if err := result.Parse(); err != nil {
				logrus.WithError(err).Error("Failed to parse message")
				continue
			}

			if ok := result.IsValid(); !ok {
				logrus.Error("Invalid message")
				continue
			}

			logrus.Infof("Got message for game %s from player %s", gameID, playerID)
			g.Coordinator.HandlePlayerMessage(gameID, playerID, result)

		// If the game is over, stop
		case <-finish:
			break
		}
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
//	@Failure	400	"Invalid websocket headers"
//	@Failure	403	"Not your game"
//	@Failure	404	"Game not found"
//	@Failure	404	"Game is not started"
//	@Failure	500	"Internal Server Error"
//	@Router		/api/v1/games/{id}/connection [get]
//	@Security	JWT
func (g *GameConnectionHandler) GetCreator(c *gin.Context) {
	authID := c.GetString("user")
	creatorID := uuid.MustParse(authID)

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

	if game.Quiz.CreatorID != creatorID {
		logrus.Error("Not your game")
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	creator, err := g.CreatorService.GetByID(creatorID)
	if err != nil {
		logrus.WithError(err).Error("How even?")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// creatorUpgrader allows us to put the subprotocol as the authentication header, it's required
	// because the browser verifies the response
	creatorUpgrader := websocket.Upgrader{
		Subprotocols: []string{c.GetHeader("Sec-Websocket-Protocol")},
		CheckOrigin: func(*http.Request) bool {
			return true
		},
	}

	ws, err := creatorUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logrus.WithError(err).Error("Game can not be joined")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	finish := make(chan struct{})

	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("Had to recover from a panic: %v", err)
		}

		g.Coordinator.UnsubscribeCreator(gameID)
	}()

	logrus.Infof("Opening websocket for creator %s in game %s", authID, gameID)
	g.Coordinator.SubscribeCreator(gameID, creator, func(message *coordinator.BroadcastMessage) {
		if err := ws.WriteJSON(message); err != nil {
			logrus.WithError(err).Error("Failed to write JSON")
		}

		// Quit the game if we're finished
		if message.Type == coordinator.FinishGameType {
			finish <- struct{}{}
		}
	})

	for {
		select {
		default:
			var result *coordinator.CreatorMessage
			if err := ws.ReadJSON(&result); err != nil {
				continue
			}

			logrus.Infof("Got message for game %s from creator %s", gameID, authID)
			g.Coordinator.HandleCreatorMessage(gameID, result)

		// If the game is over, stop
		case <-finish:
			break
		}
	}
}
