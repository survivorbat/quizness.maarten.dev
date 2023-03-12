package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGame_PlayerJoin_AddsPlayer(t *testing.T) {
	t.Parallel()
	// Arrange
	player := &Player{Nickname: "abc"}
	game := Game{}

	// Act
	err := game.PlayerJoin(player)

	// Assert
	assert.NoError(t, err)
}

func TestGame_PlayerJoin_ReturnsErrorOnAlreadyPresent(t *testing.T) {
	t.Parallel()
	// Arrange
	player := &Player{Nickname: "abc"}
	game := Game{Players: []*Player{{Nickname: "abc"}}}

	// Act
	err := game.PlayerJoin(player)

	// Assert
	assert.ErrorContains(t, err, "player is already in this game")
}

func TestGame_PlayerLeave_AddsPlayer(t *testing.T) {
	t.Parallel()
	// Arrange
	player := &Player{Nickname: "abc"}
	game := Game{Players: []*Player{{Nickname: "abc"}}}

	// Act
	err := game.PlayerLeave(player)

	// Assert
	assert.NoError(t, err)
}

func TestGame_PlayerLeave_ReturnsErrorOnAlreadyPresent(t *testing.T) {
	t.Parallel()
	// Arrange
	player := &Player{Nickname: "abc"}
	game := Game{Players: []*Player{}}

	// Act
	err := game.PlayerLeave(player)

	// Assert
	assert.ErrorContains(t, err, "player is not in this game")
}

func TestGame_Start_SetsStartTimeAndCode(t *testing.T) {
	t.Parallel()
	// Arrange
	game := Game{}

	// Act
	err := game.Start()

	// Assert
	assert.NoError(t, err)
	assert.False(t, game.StartTime.IsZero())
	assert.Len(t, game.Code, 6)
}

func TestGame_Start_ErrorsOnAlreadyStarted(t *testing.T) {
	t.Parallel()
	// Arrange
	game := Game{}

	_ = game.Start()

	// Act
	err := game.Start()

	// Assert
	assert.ErrorContains(t, err, "game has already started")
}

func TestGame_Finish_SetsFinishTime(t *testing.T) {
	t.Parallel()
	// Arrange
	game := Game{StartTime: time.Now()}

	// Act
	err := game.Finish()

	// Assert
	assert.NoError(t, err)
	assert.False(t, game.FinishTime.IsZero())
}

func TestGame_Finish_ErrorsOnNotStarted(t *testing.T) {
	t.Parallel()
	// Arrange
	game := Game{}

	// Act
	err := game.Finish()

	// Assert
	assert.ErrorContains(t, err, "game has not started")
}

func TestGame_Finish_ErrorsOnAlreadyFinished(t *testing.T) {
	t.Parallel()
	// Arrange
	game := Game{StartTime: time.Now()}

	_ = game.Finish()

	// Act
	err := game.Finish()

	// Assert
	assert.ErrorContains(t, err, "game has already finished")
}
