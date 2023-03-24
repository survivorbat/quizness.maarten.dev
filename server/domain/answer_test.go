package domain

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGameAnswers_Contains_ReturnsTrueOnContains(t *testing.T) {
	t.Parallel()
	// Arrange
	questionID := uuid.MustParse("d9d08d5c-0cab-4f1a-9b46-2bdbc10cfa93")
	playerID := uuid.MustParse("378d5245-a847-441b-9895-6351f92c2148")

	answers := GameAnswers{
		{QuestionID: questionID, PlayerID: playerID},
	}

	// Act
	result := answers.Contains(questionID, playerID)

	// Assert
	assert.True(t, result)
}

func TestGameAnswers_Contains_ReturnsFalseOnNotContains(t *testing.T) {
	t.Parallel()
	// Arrange
	questionID := uuid.MustParse("d9d08d5c-0cab-4f1a-9b46-2bdbc10cfa93")
	playerID := uuid.MustParse("378d5245-a847-441b-9895-6351f92c2148")

	answers := GameAnswers{}

	// Act
	result := answers.Contains(questionID, playerID)

	// Assert
	assert.False(t, result)
}
