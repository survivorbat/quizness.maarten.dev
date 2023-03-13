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

func TestQuiz_GetQuestion_ReturnsExpectedValue(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		quiz     *Quiz
		index    uint
		expected string
	}{
		"first": {
			quiz: &Quiz{MultipleChoiceQuestions: []*MultipleChoiceQuestion{
				{BaseQuestion: BaseQuestion{Title: "First question!", Order: 0}},
			}},
			expected: "First question!",
			index:    0,
		},
		"second": {
			quiz: &Quiz{MultipleChoiceQuestions: []*MultipleChoiceQuestion{
				{BaseQuestion: BaseQuestion{Title: "Second question!", Order: 1}},
				{BaseQuestion: BaseQuestion{Title: "First question!", Order: 0}},
				{BaseQuestion: BaseQuestion{Title: "Third question!", Order: 2}},
			}},
			expected: "Second question!",
			index:    1,
		},
	}

	for name, testData := range tests {
		testData := testData
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			// Act
			result, err := testData.quiz.GetQuestion(testData.index)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, testData.expected, result.GetBaseQuestion().Title)
		})
	}
}
