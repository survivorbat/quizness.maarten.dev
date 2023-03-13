package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGame_IsInProgress_ReturnsTrueOnStarted(t *testing.T) {
	t.Parallel()
	// Arrange
	game := &Game{StartTime: time.Now()}

	// Act
	result := game.IsInProgress()

	// Assert
	assert.True(t, result)
}

func TestGame_IsInProgress_ReturnsFalseOnNotStarted(t *testing.T) {
	t.Parallel()
	// Arrange
	game := &Game{}

	// Act
	result := game.IsInProgress()

	// Assert
	assert.False(t, result)
}
func TestGame_IsInProgress_ReturnsFalseOnStartedAndFinished(t *testing.T) {
	t.Parallel()
	// Arrange
	game := &Game{StartTime: time.Now(), FinishTime: time.Now()}

	// Act
	result := game.IsInProgress()

	// Assert
	assert.False(t, result)
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
