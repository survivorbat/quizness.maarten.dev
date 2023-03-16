package coordinator

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/survivorbat/go-tsyncmap"
	"github.com/survivorbat/qq.maarten.dev/server/services"
)

type BroadcastCallback func(*BroadcastMessage)

type GameCoordinator interface {
	SubscribePlayer(gameID uuid.UUID, playerID uuid.UUID, callback BroadcastCallback)
	SubscribeCreator(gameID uuid.UUID, callback BroadcastCallback)

	HandlePlayerMessage(game uuid.UUID, player uuid.UUID, message *PlayerMessage)
	HandleCreatorMessage(game uuid.UUID, message *CreatorMessage)
}

type LocalGameCoordinator struct {
	GameService services.GameService
	creators    tsyncmap.Map[uuid.UUID, BroadcastCallback]
	clients     tsyncmap.Map[uuid.UUID, *tsyncmap.Map[uuid.UUID, BroadcastCallback]]
}

func (c *LocalGameCoordinator) SubscribePlayer(gameID uuid.UUID, playerID uuid.UUID, callback BroadcastCallback) {
	value, _ := c.clients.LoadOrStore(gameID, &tsyncmap.Map[uuid.UUID, BroadcastCallback]{})
	value.Store(playerID, callback)
}

func (c *LocalGameCoordinator) SubscribeCreator(gameID uuid.UUID, callback BroadcastCallback) {
	c.creators.Store(gameID, callback)
}

func (c *LocalGameCoordinator) HandleCreatorMessage(gameID uuid.UUID, message *CreatorMessage) {
	game, err := c.GameService.GetByID(gameID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get game")
		return
	}

	switch message.Action {
	case NextQuestionAction:
		if err := c.GameService.Next(game); err != nil {
			logrus.WithError(err).Error("Failed to answer question")
			return
		}

		broadcast := &BroadcastMessage{
			Type:       NextQuestionType,
			QuestionID: game.CurrentQuestion,
		}

		if err := c.broadcast(gameID, broadcast); err != nil {
			logrus.WithError(err).Error("Failed to broadcast")
		}
	}
}

func (c *LocalGameCoordinator) HandlePlayerMessage(gameID uuid.UUID, playerID uuid.UUID, message *PlayerMessage) {
	game, err := c.GameService.GetByID(gameID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get game")
		return
	}

	switch message.Action {
	case AnswerAction:
		if err := c.GameService.AnswerQuestion(game, game.CurrentQuestion, playerID, message.Answer.OptionID); err != nil {
			logrus.WithError(err).Error("Failed to answer question")
			return
		}

		fmt.Println("yes can do")

		broadcast := &BroadcastMessage{
			Type:     PlayerAnsweredType,
			PlayerID: playerID,
		}

		if err := c.broadcast(gameID, broadcast); err != nil {
			logrus.WithError(err).Error("Failed to broadcast")
		}
	}
}

func (c *LocalGameCoordinator) broadcast(game uuid.UUID, message *BroadcastMessage) error {
	result, ok := c.clients.Load(game)
	if !ok {
		err := errors.New("game not found")
		logrus.WithError(err).Error("Game not found in clients, somehow")
		return err
	}

	result.Range(func(_ uuid.UUID, callback BroadcastCallback) bool {
		callback(message)
		return true
	})

	creator, ok := c.creators.Load(game)
	if !ok {
		err := errors.New("game not found")
		logrus.WithError(err).Error("Game not found in creators, somehow")
		return err
	}

	creator(message)
	return nil
}
