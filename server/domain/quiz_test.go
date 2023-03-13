package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestQuiz_HasGameInProgress_ReturnsTrueOnInProgress(t *testing.T) {
	t.Parallel()
	// Arrange
	quiz := &Quiz{Games: []*Game{{}, {StartTime: time.Now()}}}

	// Act
	result := quiz.HasGameInProgress()

	// Assert
	assert.True(t, result)
}

func TestQuiz_HasGameInProgress_ReturnsFalseOnNotInProgress(t *testing.T) {
	t.Parallel()
	// Arrange
	quiz := &Quiz{Games: []*Game{{}, {}}}

	// Act
	result := quiz.HasGameInProgress()

	// Assert
	assert.False(t, result)
}
