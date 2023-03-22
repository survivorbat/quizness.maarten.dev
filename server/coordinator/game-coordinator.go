package coordinator

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/survivorbat/go-tsyncmap"
	"github.com/survivorbat/qq.maarten.dev/server/domain"
	"github.com/survivorbat/qq.maarten.dev/server/services"
)

// Compile-time interface checks
var _ GameCoordinator = new(LocalGameCoordinator)

type BroadcastCallback func(*BroadcastMessage)

type GameCoordinator interface {
	SubscribePlayer(gameID uuid.UUID, player *domain.Player, callback BroadcastCallback)
	SubscribeCreator(gameID uuid.UUID, creator *domain.Creator, callback BroadcastCallback)

	UnsubscribePlayer(gameID uuid.UUID, player *domain.Player)
	UnsubscribeCreator(gameID uuid.UUID)

	HandlePlayerMessage(game uuid.UUID, player uuid.UUID, message *PlayerMessage)
	HandleCreatorMessage(game uuid.UUID, message *CreatorMessage)
}

// creatorInfo is a container for the creator and their callback
type creatorInfo struct {
	callback BroadcastCallback
	creator  *domain.Creator
}

// LocalGameCoordinator coordinates a running game
type LocalGameCoordinator struct {
	// GameService is used to manipulate games
	GameService services.GameService

	// creators is a mapping of game ids with creatorInfo
	creators tsyncmap.Map[uuid.UUID, creatorInfo]

	// clients is a list of games with connected players and callbacks
	clients tsyncmap.Map[uuid.UUID, *tsyncmap.Map[*domain.Player, BroadcastCallback]]
}

func (c *LocalGameCoordinator) SubscribePlayer(gameID uuid.UUID, player *domain.Player, callback BroadcastCallback) {
	value, _ := c.clients.LoadOrStore(gameID, &tsyncmap.Map[*domain.Player, BroadcastCallback]{})
	value.Store(player, callback)
	c.broadcastState(gameID)
}

func (c *LocalGameCoordinator) SubscribeCreator(gameID uuid.UUID, creator *domain.Creator, callback BroadcastCallback) {
	c.creators.Store(gameID, creatorInfo{creator: creator, callback: callback})
	c.broadcastState(gameID)
}

func (c *LocalGameCoordinator) UnsubscribePlayer(gameID uuid.UUID, player *domain.Player) {
	value, ok := c.clients.Load(gameID)
	if ok {
		value.Delete(player)
	}

	c.broadcastState(gameID)
}

func (c *LocalGameCoordinator) UnsubscribeCreator(gameID uuid.UUID) {
	c.creators.Delete(gameID)
	c.broadcastState(gameID)
}

func (c *LocalGameCoordinator) HandleCreatorMessage(gameID uuid.UUID, message *CreatorMessage) {
	game, err := c.GameService.GetByID(gameID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get game")
		return
	}

	switch message.Action {
	case FinishGameAction:
		if err := c.GameService.Finish(game); err != nil {
			logrus.WithError(err).Error("Failed to finish")
			return
		}

		broadcast := &BroadcastMessage{
			Type: FinishGameType,
		}

		c.broadcast(gameID, broadcast)

	case NextQuestionAction:
		if err := c.GameService.Next(game); err != nil {
			logrus.WithError(err).Error("Failed to answer question")
			return
		}

		// Broadcast the new state
		c.broadcastState(game.ID)
	}
}

func (c *LocalGameCoordinator) HandlePlayerMessage(gameID uuid.UUID, player uuid.UUID, message *PlayerMessage) {
	game, err := c.GameService.GetByID(gameID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get game")
		return
	}

	switch message.Action {
	case AnswerAction:
		if err := c.GameService.AnswerQuestion(game, game.CurrentQuestion, player, message.Answer.OptionID); err != nil {
			logrus.WithError(err).Error("Failed to answer question")
			return
		}

		broadcast := &BroadcastMessage{
			Type:                  PlayerAnsweredType,
			PlayerAnsweredContent: &playerAnsweredContent{PlayerID: player},
		}

		c.broadcast(gameID, broadcast)
	}
}

func (c *LocalGameCoordinator) broadcastState(gameID uuid.UUID) {
	game, err := c.GameService.GetByID(gameID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get game")
		return
	}

	message := &BroadcastMessage{
		Type: StateType,
		StateContent: &stateContent{
			Players:         []*participant{},
			CurrentQuestion: game.CurrentQuestion,
			CurrentDeadline: game.CurrentDeadline,
		},
	}

	creator, ok := c.creators.Load(gameID)
	if ok {
		message.StateContent.Creator = &participant{
			ID:              creator.creator.ID,
			Nickname:        creator.creator.Nickname,
			Color:           creator.creator.Color,
			BackgroundColor: creator.creator.BackgroundColor,
		}
	}

	result, ok := c.clients.Load(gameID)
	if ok {
		result.Range(func(player *domain.Player, broadcast BroadcastCallback) bool {
			message.StateContent.Players = append(message.StateContent.Players, &participant{
				ID:              player.ID,
				Nickname:        player.Nickname,
				Color:           player.Color,
				BackgroundColor: player.BackgroundColor,
			})
			return true
		})
	}

	c.broadcast(gameID, message)
}

// broadcast sends a message to the creator and
func (c *LocalGameCoordinator) broadcast(game uuid.UUID, message *BroadcastMessage) {
	var (
		playerCount      int
		creatorAvailable bool
	)

	result, ok := c.clients.Load(game)
	if ok {
		result.Range(func(player *domain.Player, broadcast BroadcastCallback) bool {
			broadcast(message)
			playerCount++
			return true
		})
	}

	creator, ok := c.creators.Load(game)
	if ok {
		creator.callback(message)
		creatorAvailable = true
	}

	logrus.Infof("Broadcast to %d players and creator (%t): %#v", playerCount, creatorAvailable, message)
}
