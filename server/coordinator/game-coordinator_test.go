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
	creatorCalledWith []*BroadcastMessage
	playerCalledWith  []*BroadcastMessage
}

func (c *callbackCollection) player(msg *BroadcastMessage) {
	c.playerCalledWith = append(c.playerCalledWith, msg)
}

func (c *callbackCollection) creator(msg *BroadcastMessage) {
	c.creatorCalledWith = append(c.creatorCalledWith, msg)
}

func TestLocalGameCoordinator_SubscribePlayer_AddsClientAndBroadcasts(t *testing.T) {
	t.Parallel()
	// Arrange
	gameID := uuid.MustParse("2389b70a-74df-439c-8d5f-cf4f3f9471bd")
	playerID := uuid.MustParse("ffcdf7eb-0eee-411f-9b3f-2401315cc9e6")

	game := &domain.Game{
		BaseObject:      domain.BaseObject{ID: gameID},
		CurrentQuestion: uuid.MustParse("ffcdf7eb-0eee-411f-9b3f-2401315cc9e6"),
		CurrentDeadline: time.Now().Add(5 * time.Hour),
	}

	coordinator := &LocalGameCoordinator{
		GameService: &MockGameService{
			getByIDReturns: game,
		},
	}
	callbacks := new(callbackCollection)

	player := &domain.Player{
		BaseObject: domain.BaseObject{ID: playerID},
		Nickname:   "Test",
	}

	// Act
	coordinator.SubscribePlayer(gameID, player, callbacks.player)

	// Assert
	item, ok := coordinator.clients.Load(gameID)
	if !assert.True(t, ok) {
		t.Fatal("not found")
	}

	result, ok := item.Load(player)
	if !assert.True(t, ok) {
		t.Fatal("not found")
	}

	assert.NotNil(t, result)

	if assert.Len(t, callbacks.playerCalledWith, 1) {
		assert.Equal(t, player.ID, callbacks.playerCalledWith[0].StateContent.Players[0].ID)
		assert.Equal(t, player.Nickname, callbacks.playerCalledWith[0].StateContent.Players[0].Nickname)
		assert.Equal(t, game.CurrentQuestion, callbacks.playerCalledWith[0].StateContent.CurrentQuestion)
		assert.Equal(t, game.CurrentDeadline, callbacks.playerCalledWith[0].StateContent.CurrentDeadline)
	}

}
func TestLocalGameCoordinator_UnsubscribePlayer_RemovesClientAndBroadcasts(t *testing.T) {
	t.Parallel()
	// Arrange
	gameID := uuid.MustParse("2389b70a-74df-439c-8d5f-cf4f3f9471bd")
	player1ID := uuid.MustParse("ffcdf7eb-0eee-411f-9b3f-2401315cc9e6")
	player2ID := uuid.MustParse("14eeb8e1-4db9-4177-bf47-ae80a22b2ff9")

	game := &domain.Game{
		BaseObject:      domain.BaseObject{ID: gameID},
		CurrentQuestion: uuid.MustParse("ffcdf7eb-0eee-411f-9b3f-2401315cc9e6"),
		CurrentDeadline: time.Now().Add(5 * time.Hour),
	}

	coordinator := &LocalGameCoordinator{
		GameService: &MockGameService{
			getByIDReturns: game,
		},
	}
	callbacks := new(callbackCollection)

	player1 := &domain.Player{
		BaseObject: domain.BaseObject{ID: player1ID},
		Nickname:   "def",
	}
	player2 := &domain.Player{
		BaseObject: domain.BaseObject{ID: player2ID},
		Nickname:   "abc",
	}

	coordinator.SubscribePlayer(gameID, player1, callbacks.player)
	coordinator.SubscribePlayer(gameID, player2, callbacks.player)

	// Act
	coordinator.UnsubscribePlayer(gameID, player2)

	// Assert
	item, ok := coordinator.clients.Load(gameID)
	if !assert.True(t, ok) {
		t.Fatal("not found")
	}

	_, ok = item.Load(player2)
	assert.False(t, ok)

	if assert.Len(t, callbacks.playerCalledWith, 4) {
		// Only one player should be left
		assert.Len(t, callbacks.playerCalledWith[3].StateContent.Players, 1)
	}
}

func TestLocalGameCoordinator_SubscribeCreator_AddsClient(t *testing.T) {
	t.Parallel()
	// Arrange
	creatorID := uuid.MustParse("f8f9cf51-31d7-4a6b-a1fc-63f5e750e16a")
	gameID := uuid.MustParse("2389b70a-74df-439c-8d5f-cf4f3f9471bd")

	game := &domain.Game{
		BaseObject:      domain.BaseObject{ID: gameID},
		CurrentQuestion: uuid.MustParse("ffcdf7eb-0eee-411f-9b3f-2401315cc9e6"),
		CurrentDeadline: time.Now().Add(5 * time.Hour),
	}

	coordinator := &LocalGameCoordinator{
		GameService: &MockGameService{
			getByIDReturns: game,
		},
	}
	callbacks := new(callbackCollection)

	creator := &domain.Creator{
		BaseObject: domain.BaseObject{ID: creatorID},
		Nickname:   "Abc",
	}

	// Act
	coordinator.SubscribeCreator(gameID, creator, callbacks.creator)

	// Assert
	item, ok := coordinator.creators.Load(gameID)
	if !assert.True(t, ok) {
		t.Fatal("not found")
	}

	assert.NotNil(t, item)

	if assert.Len(t, callbacks.creatorCalledWith, 1) {
		assert.Equal(t, creator.ID, callbacks.creatorCalledWith[0].StateContent.Creator.ID)
		assert.Equal(t, creator.Nickname, callbacks.creatorCalledWith[0].StateContent.Creator.Nickname)
	}
}

func TestLocalGameCoordinator_UnsubscribeCreator_RemovesClient(t *testing.T) {
	t.Parallel()
	// Arrange
	creatorID := uuid.MustParse("f8f9cf51-31d7-4a6b-a1fc-63f5e750e16a")
	gameID := uuid.MustParse("2389b70a-74df-439c-8d5f-cf4f3f9471bd")
	playerID := uuid.MustParse("4b629c20-9100-46c8-a08c-69bcd19e7043")

	coordinator := &LocalGameCoordinator{
		GameService: &MockGameService{
			getByIDReturns: &domain.Game{BaseObject: domain.BaseObject{ID: gameID}},
		},
	}
	callbacks := new(callbackCollection)

	creator := &domain.Creator{
		BaseObject: domain.BaseObject{ID: creatorID},
		Nickname:   "Abc",
	}

	player := &domain.Player{
		BaseObject: domain.BaseObject{ID: playerID},
		Nickname:   "Test",
	}

	coordinator.SubscribeCreator(gameID, creator, callbacks.creator)
	coordinator.SubscribePlayer(gameID, player, callbacks.player)

	// Act
	coordinator.UnsubscribeCreator(gameID)

	// Assert
	_, ok := coordinator.creators.Load(gameID)
	assert.False(t, ok)

	if assert.Len(t, callbacks.playerCalledWith, 2) {
		// Should no longer contain the creator
		assert.Nil(t, callbacks.playerCalledWith[1].StateContent.Creator)
	}
}

func TestLocalGameCoordinator_HandlePlayerMessage_AnswerLaunchesBroadcast(t *testing.T) {
	t.Parallel()
	// Arrange
	gameID := uuid.MustParse("2389b70a-74df-439c-8d5f-cf4f3f9471bd")
	playerID := uuid.MustParse("ffcdf7eb-0eee-411f-9b3f-2401315cc9e6")
	questionID := uuid.MustParse("67ec56fa-d082-4fcd-b373-885801e7a910")

	player := &domain.Player{
		BaseObject: domain.BaseObject{ID: playerID},
	}

	game := &domain.Game{
		BaseObject:      domain.BaseObject{ID: gameID},
		CurrentQuestion: questionID,
	}

	gameService := &MockGameService{getByIDReturns: game}
	coordinator := &LocalGameCoordinator{GameService: gameService}
	callbacks := new(callbackCollection)

	coordinator.SubscribeCreator(gameID, &domain.Creator{}, callbacks.creator)
	coordinator.SubscribePlayer(gameID, player, callbacks.player)

	message := &PlayerMessage{
		Action: AnswerAction,
		Answer: &inputs.Answer{
			OptionID: uuid.MustParse("37d3d2e9-0466-4eb6-90ae-a9f63d036de5"),
		},
	}

	// Act
	coordinator.HandlePlayerMessage(gameID, player.ID, message)

	// Assert
	assert.Equal(t, game, gameService.answerQuestionCalledWithGame)
	assert.Equal(t, player.ID, gameService.answerQuestionCalledWithPlayer)
	assert.Equal(t, questionID, gameService.answerQuestionCalledWithQuestion)
	assert.Equal(t, message.Answer.OptionID, gameService.answerQuestionCalledWithOption)

	if assert.Len(t, callbacks.playerCalledWith, 2) {
		assert.Equal(t, player.ID, callbacks.playerCalledWith[1].PlayerAnsweredContent.PlayerID)
		assert.Equal(t, PlayerAnsweredType, callbacks.playerCalledWith[1].Type)
	}
	if assert.Len(t, callbacks.creatorCalledWith, 3) {
		assert.Equal(t, player.ID, callbacks.creatorCalledWith[2].PlayerAnsweredContent.PlayerID)
		assert.Equal(t, PlayerAnsweredType, callbacks.creatorCalledWith[2].Type)
	}
}

func TestLocalGameCoordinator_HandlePlayerMessage_DoesNothingOnGameNotFound(t *testing.T) {
	t.Parallel()
	// Arrange
	gameID := uuid.MustParse("2389b70a-74df-439c-8d5f-cf4f3f9471bd")
	playerID := uuid.MustParse("ffcdf7eb-0eee-411f-9b3f-2401315cc9e6")

	player := &domain.Player{
		BaseObject: domain.BaseObject{ID: playerID},
	}

	gameService := &MockGameService{getByIDReturnsError: assert.AnError}

	coordinator := &LocalGameCoordinator{GameService: gameService}
	callbacks := new(callbackCollection)

	coordinator.SubscribeCreator(gameID, &domain.Creator{}, callbacks.creator)
	coordinator.SubscribePlayer(gameID, player, callbacks.player)

	message := &PlayerMessage{
		Action: AnswerAction,
		Answer: &inputs.Answer{
			OptionID: uuid.MustParse("37d3d2e9-0466-4eb6-90ae-a9f63d036de5"),
		},
	}

	// Act
	coordinator.HandlePlayerMessage(gameID, player.ID, message)

	// Assert
	assert.Empty(t, gameService.answerQuestionCalledWithGame)
	assert.Empty(t, gameService.answerQuestionCalledWithPlayer)
	assert.Empty(t, gameService.answerQuestionCalledWithQuestion)
	assert.Empty(t, gameService.answerQuestionCalledWithOption)

	assert.Len(t, callbacks.playerCalledWith, 0)
	assert.Len(t, callbacks.creatorCalledWith, 0)
}

func TestLocalGameCoordinator_HandlePlayerMessage_DoesNothingOnAnswerFailure(t *testing.T) {
	t.Parallel()
	// Arrange
	gameID := uuid.MustParse("2389b70a-74df-439c-8d5f-cf4f3f9471bd")
	playerID := uuid.MustParse("ffcdf7eb-0eee-411f-9b3f-2401315cc9e6")
	questionID := uuid.MustParse("67ec56fa-d082-4fcd-b373-885801e7a910")

	player := &domain.Player{
		BaseObject: domain.BaseObject{ID: playerID},
	}

	game := &domain.Game{
		BaseObject:      domain.BaseObject{ID: gameID},
		CurrentQuestion: questionID,
	}

	gameService := &MockGameService{getByIDReturns: game, answerQuestionReturns: assert.AnError}

	coordinator := &LocalGameCoordinator{GameService: gameService}
	callbacks := new(callbackCollection)

	coordinator.SubscribeCreator(gameID, &domain.Creator{}, callbacks.creator)
	coordinator.SubscribePlayer(gameID, player, callbacks.player)

	message := &PlayerMessage{
		Action: AnswerAction,
		Answer: &inputs.Answer{
			OptionID: uuid.MustParse("37d3d2e9-0466-4eb6-90ae-a9f63d036de5"),
		},
	}

	// Act
	coordinator.HandlePlayerMessage(gameID, player.ID, message)

	// Assert
	assert.Equal(t, game, gameService.answerQuestionCalledWithGame)
	assert.Equal(t, player.ID, gameService.answerQuestionCalledWithPlayer)
	assert.Equal(t, questionID, gameService.answerQuestionCalledWithQuestion)
	assert.Equal(t, message.Answer.OptionID, gameService.answerQuestionCalledWithOption)

	assert.Len(t, callbacks.playerCalledWith, 1)
	assert.Len(t, callbacks.creatorCalledWith, 2)
}

func TestLocalGameCoordinator_HandleCreatorMessage_NextLaunchesBroadcast(t *testing.T) {
	t.Parallel()
	// Arrange
	gameID := uuid.MustParse("2389b70a-74df-439c-8d5f-cf4f3f9471bd")
	questionID := uuid.MustParse("67ec56fa-d082-4fcd-b373-885801e7a910")
	playerID := uuid.MustParse("ffcdf7eb-0eee-411f-9b3f-2401315cc9e6")

	player := &domain.Player{
		BaseObject: domain.BaseObject{ID: playerID},
	}

	game := &domain.Game{
		BaseObject: domain.BaseObject{ID: gameID},
	}

	gameService := &MockGameService{getByIDReturns: game, nextSetsCurrentQuestion: questionID}
	coordinator := &LocalGameCoordinator{GameService: gameService}
	callbacks := new(callbackCollection)

	coordinator.SubscribeCreator(gameID, &domain.Creator{}, callbacks.creator)
	coordinator.SubscribePlayer(gameID, player, callbacks.player)

	message := &CreatorMessage{
		Action: NextQuestionAction,
	}

	// Act
	coordinator.HandleCreatorMessage(gameID, message)

	// Assert
	if assert.Len(t, callbacks.playerCalledWith, 2) {
		assert.Equal(t, questionID, callbacks.playerCalledWith[1].NextQuestionContent.QuestionID)
		assert.Equal(t, NextQuestionType, callbacks.playerCalledWith[1].Type)
	}

	if assert.Len(t, callbacks.creatorCalledWith, 3) {
		assert.Equal(t, questionID, callbacks.creatorCalledWith[2].NextQuestionContent.QuestionID)
		assert.Equal(t, NextQuestionType, callbacks.creatorCalledWith[2].Type)
	}
}

func TestLocalGameCoordinator_HandleCreatorMessage_DoesNothingOnGameNotFound(t *testing.T) {
	t.Parallel()
	// Arrange
	gameID := uuid.MustParse("2389b70a-74df-439c-8d5f-cf4f3f9471bd")
	playerID := uuid.MustParse("ffcdf7eb-0eee-411f-9b3f-2401315cc9e6")

	player := &domain.Player{
		BaseObject: domain.BaseObject{ID: playerID},
	}

	gameService := &MockGameService{getByIDReturnsError: assert.AnError}

	coordinator := &LocalGameCoordinator{GameService: gameService}
	callbacks := new(callbackCollection)

	coordinator.SubscribeCreator(gameID, &domain.Creator{}, callbacks.creator)
	coordinator.SubscribePlayer(gameID, player, callbacks.player)

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

	assert.Len(t, callbacks.playerCalledWith, 0)
	assert.Len(t, callbacks.creatorCalledWith, 0)
}

func TestLocalGameCoordinator_HandleCreatorMessage_DoesNothingOnNextFailure(t *testing.T) {
	t.Parallel()
	// Arrange
	gameID := uuid.MustParse("2389b70a-74df-439c-8d5f-cf4f3f9471bd")
	playerID := uuid.MustParse("ffcdf7eb-0eee-411f-9b3f-2401315cc9e6")
	questionID := uuid.MustParse("67ec56fa-d082-4fcd-b373-885801e7a910")

	player := &domain.Player{
		BaseObject: domain.BaseObject{ID: playerID},
	}

	game := &domain.Game{
		BaseObject:      domain.BaseObject{ID: gameID},
		CurrentQuestion: questionID,
	}

	gameService := &MockGameService{getByIDReturns: game, nextReturns: assert.AnError}

	coordinator := &LocalGameCoordinator{GameService: gameService}
	callbacks := new(callbackCollection)

	coordinator.SubscribeCreator(gameID, &domain.Creator{}, callbacks.creator)
	coordinator.SubscribePlayer(gameID, player, callbacks.player)

	message := &CreatorMessage{
		Action: NextQuestionAction,
	}

	// Act
	coordinator.HandleCreatorMessage(gameID, message)

	// Assert
	assert.Len(t, callbacks.playerCalledWith, 1)
	assert.Len(t, callbacks.creatorCalledWith, 2)
}

func TestLocalGameCoordinator_HandleCreatorMessage_FinishLaunchesBroadcast(t *testing.T) {
	t.Parallel()
	// Arrange
	gameID := uuid.MustParse("2389b70a-74df-439c-8d5f-cf4f3f9471bd")
	playerID := uuid.MustParse("ffcdf7eb-0eee-411f-9b3f-2401315cc9e6")

	player := &domain.Player{
		BaseObject: domain.BaseObject{ID: playerID},
	}

	game := &domain.Game{
		BaseObject: domain.BaseObject{ID: gameID},
		StartTime:  time.Now(),
	}

	gameService := &MockGameService{getByIDReturns: game}
	coordinator := &LocalGameCoordinator{GameService: gameService}
	callbacks := new(callbackCollection)

	coordinator.SubscribeCreator(gameID, &domain.Creator{}, callbacks.creator)
	coordinator.SubscribePlayer(gameID, player, callbacks.player)

	message := &CreatorMessage{
		Action: FinishGameAction,
	}

	// Act
	coordinator.HandleCreatorMessage(gameID, message)

	// Assert
	if assert.Len(t, callbacks.playerCalledWith, 2) {
		assert.Equal(t, FinishGameType, callbacks.playerCalledWith[1].Type)
	}

	if assert.Len(t, callbacks.creatorCalledWith, 3) {
		assert.Equal(t, FinishGameType, callbacks.creatorCalledWith[2].Type)
	}
}

func TestLocalGameCoordinator_HandleCreatorMessage_DoesNothingOnFinishFailure(t *testing.T) {
	t.Parallel()
	// Arrange
	gameID := uuid.MustParse("2389b70a-74df-439c-8d5f-cf4f3f9471bd")
	playerID := uuid.MustParse("ffcdf7eb-0eee-411f-9b3f-2401315cc9e6")

	player := &domain.Player{
		BaseObject: domain.BaseObject{ID: playerID},
	}

	game := &domain.Game{
		BaseObject:      domain.BaseObject{ID: gameID},
		CurrentDeadline: time.Now().Add(time.Hour * 5),
	}

	gameService := &MockGameService{getByIDReturns: game, finishReturns: assert.AnError}

	coordinator := &LocalGameCoordinator{GameService: gameService}
	callbacks := new(callbackCollection)

	coordinator.SubscribeCreator(gameID, &domain.Creator{}, callbacks.creator)
	coordinator.SubscribePlayer(gameID, player, callbacks.player)

	message := &CreatorMessage{
		Action: FinishGameAction,
	}

	// Act
	coordinator.HandleCreatorMessage(gameID, message)

	// Assert
	assert.Len(t, callbacks.playerCalledWith, 1)
	assert.Len(t, callbacks.creatorCalledWith, 2)
}
