package coordinator

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/survivorbat/qq.maarten.dev/server/domain"
	"github.com/survivorbat/qq.maarten.dev/server/inputs"
	"testing"
	"time"
)

type callbackCollection struct {
	creatorCalledWith *BroadcastMessage
	playerCalledWith  *BroadcastMessage
}

func (c *callbackCollection) player(msg *BroadcastMessage) {
	c.playerCalledWith = msg
}

func (c *callbackCollection) creator(msg *BroadcastMessage) {
	c.creatorCalledWith = msg
}

func TestLocalGameCoordinator_SubscribePlayer_AddsClient(t *testing.T) {
	t.Parallel()
	// Arrange
	gameID := uuid.MustParse("2389b70a-74df-439c-8d5f-cf4f3f9471bd")
	playerID := uuid.MustParse("ffcdf7eb-0eee-411f-9b3f-2401315cc9e6")

	coordinator := &LocalGameCoordinator{}
	callbacks := new(callbackCollection)

	// Act
	coordinator.SubscribePlayer(gameID, playerID, callbacks.player)

	// Assert
	item, ok := coordinator.clients.Load(gameID)
	if !assert.True(t, ok) {
		t.Fatal("not found")
	}

	result, ok := item.Load(playerID)
	if !assert.True(t, ok) {
		t.Fatal("not found")
	}

	assert.NotNil(t, result)
}

func TestLocalGameCoordinator_SubscribeCreator_AddsClient(t *testing.T) {
	t.Parallel()
	// Arrange
	gameID := uuid.MustParse("2389b70a-74df-439c-8d5f-cf4f3f9471bd")

	coordinator := &LocalGameCoordinator{}
	callbacks := new(callbackCollection)

	// Act
	coordinator.SubscribeCreator(gameID, callbacks.creator)

	// Assert
	item, ok := coordinator.creators.Load(gameID)
	if !assert.True(t, ok) {
		t.Fatal("not found")
	}

	assert.NotNil(t, item)
}

func TestLocalGameCoordinator_HandlePlayerMessage_AnswerLaunchesBroadcast(t *testing.T) {
	t.Parallel()
	// Arrange
	gameID := uuid.MustParse("2389b70a-74df-439c-8d5f-cf4f3f9471bd")
	playerID := uuid.MustParse("ffcdf7eb-0eee-411f-9b3f-2401315cc9e6")
	questionID := uuid.MustParse("67ec56fa-d082-4fcd-b373-885801e7a910")

	game := &domain.Game{
		BaseObject:      domain.BaseObject{ID: gameID},
		CurrentQuestion: questionID,
	}

	gameService := &MockGameService{getByIDReturns: game}
	coordinator := &LocalGameCoordinator{GameService: gameService}
	callbacks := new(callbackCollection)

	coordinator.SubscribePlayer(gameID, playerID, callbacks.player)
	coordinator.SubscribeCreator(gameID, callbacks.creator)

	message := &PlayerMessage{
		Action: AnswerAction,
		Answer: &inputs.Answer{
			OptionID: uuid.MustParse("37d3d2e9-0466-4eb6-90ae-a9f63d036de5"),
		},
	}

	// Act
	coordinator.HandlePlayerMessage(gameID, playerID, message)

	// Assert
	assert.Equal(t, game, gameService.answerQuestionCalledWithGame)
	assert.Equal(t, playerID, gameService.answerQuestionCalledWithPlayer)
	assert.Equal(t, questionID, gameService.answerQuestionCalledWithQuestion)
	assert.Equal(t, message.Answer.OptionID, gameService.answerQuestionCalledWithOption)

	if assert.NotNil(t, callbacks.playerCalledWith) {
		assert.Equal(t, playerID, callbacks.playerCalledWith.PlayerID)
		assert.Equal(t, PlayerAnsweredType, callbacks.playerCalledWith.Type)
	}
	if assert.NotNil(t, callbacks.creatorCalledWith) {
		assert.Equal(t, playerID, callbacks.creatorCalledWith.PlayerID)
		assert.Equal(t, PlayerAnsweredType, callbacks.creatorCalledWith.Type)
	}
}

func TestLocalGameCoordinator_HandlePlayerMessage_DoesNothingOnGameNotFound(t *testing.T) {
	t.Parallel()
	// Arrange
	gameID := uuid.MustParse("2389b70a-74df-439c-8d5f-cf4f3f9471bd")
	playerID := uuid.MustParse("ffcdf7eb-0eee-411f-9b3f-2401315cc9e6")

	gameService := &MockGameService{getByIDReturnsError: assert.AnError}

	coordinator := &LocalGameCoordinator{GameService: gameService}
	callbacks := new(callbackCollection)

	coordinator.SubscribePlayer(gameID, playerID, callbacks.player)
	coordinator.SubscribeCreator(gameID, callbacks.creator)

	message := &PlayerMessage{
		Action: AnswerAction,
		Answer: &inputs.Answer{
			OptionID: uuid.MustParse("37d3d2e9-0466-4eb6-90ae-a9f63d036de5"),
		},
	}

	// Act
	coordinator.HandlePlayerMessage(gameID, playerID, message)

	// Assert
	assert.Empty(t, gameService.answerQuestionCalledWithGame)
	assert.Empty(t, gameService.answerQuestionCalledWithPlayer)
	assert.Empty(t, gameService.answerQuestionCalledWithQuestion)
	assert.Empty(t, gameService.answerQuestionCalledWithOption)

	assert.Nil(t, callbacks.playerCalledWith)
	assert.Nil(t, callbacks.creatorCalledWith)
}

func TestLocalGameCoordinator_HandlePlayerMessage_DoesNothingOnAnswerFailure(t *testing.T) {
	t.Parallel()
	// Arrange
	gameID := uuid.MustParse("2389b70a-74df-439c-8d5f-cf4f3f9471bd")
	playerID := uuid.MustParse("ffcdf7eb-0eee-411f-9b3f-2401315cc9e6")
	questionID := uuid.MustParse("67ec56fa-d082-4fcd-b373-885801e7a910")

	game := &domain.Game{
		BaseObject:      domain.BaseObject{ID: gameID},
		CurrentQuestion: questionID,
	}

	gameService := &MockGameService{getByIDReturns: game, answerQuestionReturns: assert.AnError}

	coordinator := &LocalGameCoordinator{GameService: gameService}
	callbacks := new(callbackCollection)

	coordinator.SubscribePlayer(gameID, playerID, callbacks.player)
	coordinator.SubscribeCreator(gameID, callbacks.creator)

	message := &PlayerMessage{
		Action: AnswerAction,
		Answer: &inputs.Answer{
			OptionID: uuid.MustParse("37d3d2e9-0466-4eb6-90ae-a9f63d036de5"),
		},
	}

	// Act
	coordinator.HandlePlayerMessage(gameID, playerID, message)

	// Assert
	assert.Equal(t, game, gameService.answerQuestionCalledWithGame)
	assert.Equal(t, playerID, gameService.answerQuestionCalledWithPlayer)
	assert.Equal(t, questionID, gameService.answerQuestionCalledWithQuestion)
	assert.Equal(t, message.Answer.OptionID, gameService.answerQuestionCalledWithOption)

	assert.Nil(t, callbacks.playerCalledWith)
	assert.Nil(t, callbacks.creatorCalledWith)
}

func TestLocalGameCoordinator_HandleCreatorMessage_NextLaunchesBroadcast(t *testing.T) {
	t.Parallel()
	// Arrange
	gameID := uuid.MustParse("2389b70a-74df-439c-8d5f-cf4f3f9471bd")
	questionID := uuid.MustParse("67ec56fa-d082-4fcd-b373-885801e7a910")
	playerID := uuid.MustParse("ffcdf7eb-0eee-411f-9b3f-2401315cc9e6")

	game := &domain.Game{
		BaseObject: domain.BaseObject{ID: gameID},
	}

	gameService := &MockGameService{getByIDReturns: game, nextSetsCurrentQuestion: questionID}
	coordinator := &LocalGameCoordinator{GameService: gameService}
	callbacks := new(callbackCollection)

	coordinator.SubscribePlayer(gameID, playerID, callbacks.player)
	coordinator.SubscribeCreator(gameID, callbacks.creator)

	message := &CreatorMessage{
		Action: NextQuestionAction,
	}

	// Act
	coordinator.HandleCreatorMessage(gameID, message)

	// Assert
	if assert.NotNil(t, callbacks.playerCalledWith) {
		assert.Equal(t, questionID, callbacks.playerCalledWith.QuestionID)
		assert.Equal(t, NextQuestionType, callbacks.playerCalledWith.Type)
	}

	if assert.NotNil(t, callbacks.creatorCalledWith) {
		assert.Equal(t, questionID, callbacks.creatorCalledWith.QuestionID)
		assert.Equal(t, NextQuestionType, callbacks.creatorCalledWith.Type)
	}
}

func TestLocalGameCoordinator_HandleCreatorMessage_DoesNothingOnGameNotFound(t *testing.T) {
	t.Parallel()
	// Arrange
	gameID := uuid.MustParse("2389b70a-74df-439c-8d5f-cf4f3f9471bd")
	playerID := uuid.MustParse("ffcdf7eb-0eee-411f-9b3f-2401315cc9e6")

	gameService := &MockGameService{getByIDReturnsError: assert.AnError}

	coordinator := &LocalGameCoordinator{GameService: gameService}
	callbacks := new(callbackCollection)

	coordinator.SubscribePlayer(gameID, playerID, callbacks.player)
	coordinator.SubscribeCreator(gameID, callbacks.creator)

	message := &CreatorMessage{
		Action: NextQuestionAction,
	}

	// Act
	coordinator.HandleCreatorMessage(gameID, message)

	// Assert
	assert.Empty(t, gameService.answerQuestionCalledWithGame)
	assert.Empty(t, gameService.answerQuestionCalledWithPlayer)
	assert.Empty(t, gameService.answerQuestionCalledWithQuestion)
	assert.Empty(t, gameService.answerQuestionCalledWithOption)

	assert.Nil(t, callbacks.creatorCalledWith)
	assert.Nil(t, callbacks.playerCalledWith)
}

func TestLocalGameCoordinator_HandleCreatorMessage_DoesNothingOnNextFailure(t *testing.T) {
	t.Parallel()
	// Arrange
	gameID := uuid.MustParse("2389b70a-74df-439c-8d5f-cf4f3f9471bd")
	playerID := uuid.MustParse("ffcdf7eb-0eee-411f-9b3f-2401315cc9e6")
	questionID := uuid.MustParse("67ec56fa-d082-4fcd-b373-885801e7a910")

	game := &domain.Game{
		BaseObject:      domain.BaseObject{ID: gameID},
		CurrentQuestion: questionID,
	}

	gameService := &MockGameService{getByIDReturns: game, nextReturns: assert.AnError}

	coordinator := &LocalGameCoordinator{GameService: gameService}
	callbacks := new(callbackCollection)

	coordinator.SubscribePlayer(gameID, playerID, callbacks.player)
	coordinator.SubscribeCreator(gameID, callbacks.creator)

	message := &CreatorMessage{
		Action: NextQuestionAction,
	}

	// Act
	coordinator.HandleCreatorMessage(gameID, message)

	// Assert
	assert.Nil(t, callbacks.creatorCalledWith)
	assert.Nil(t, callbacks.playerCalledWith)
}

func TestLocalGameCoordinator_HandleCreatorMessage_FinishLaunchesBroadcast(t *testing.T) {
	t.Parallel()
	// Arrange
	gameID := uuid.MustParse("2389b70a-74df-439c-8d5f-cf4f3f9471bd")
	playerID := uuid.MustParse("ffcdf7eb-0eee-411f-9b3f-2401315cc9e6")

	game := &domain.Game{
		BaseObject: domain.BaseObject{ID: gameID},
		StartTime:  time.Now(),
	}

	gameService := &MockGameService{getByIDReturns: game}
	coordinator := &LocalGameCoordinator{GameService: gameService}
	callbacks := new(callbackCollection)

	coordinator.SubscribePlayer(gameID, playerID, callbacks.player)
	coordinator.SubscribeCreator(gameID, callbacks.creator)

	message := &CreatorMessage{
		Action: FinishGameAction,
	}

	// Act
	coordinator.HandleCreatorMessage(gameID, message)

	// Assert
	if assert.NotNil(t, callbacks.playerCalledWith) {
		assert.Equal(t, FinishGameType, callbacks.playerCalledWith.Type)
	}

	if assert.NotNil(t, callbacks.creatorCalledWith) {
		assert.Equal(t, FinishGameType, callbacks.creatorCalledWith.Type)
	}
}

func TestLocalGameCoordinator_HandleCreatorMessage_DoesNothingOnFinishFailure(t *testing.T) {
	t.Parallel()
	// Arrange
	gameID := uuid.MustParse("2389b70a-74df-439c-8d5f-cf4f3f9471bd")
	playerID := uuid.MustParse("ffcdf7eb-0eee-411f-9b3f-2401315cc9e6")

	game := &domain.Game{
		BaseObject:      domain.BaseObject{ID: gameID},
		CurrentDeadline: time.Now().Add(time.Hour * 5),
	}

	gameService := &MockGameService{getByIDReturns: game, finishReturns: assert.AnError}

	coordinator := &LocalGameCoordinator{GameService: gameService}
	callbacks := new(callbackCollection)

	coordinator.SubscribePlayer(gameID, playerID, callbacks.player)
	coordinator.SubscribeCreator(gameID, callbacks.creator)

	message := &CreatorMessage{
		Action: FinishGameAction,
	}

	// Act
	coordinator.HandleCreatorMessage(gameID, message)

	// Assert
	assert.Nil(t, callbacks.creatorCalledWith)
	assert.Nil(t, callbacks.playerCalledWith)
}
