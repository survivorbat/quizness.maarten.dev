package domain

import (
	"github.com/google/uuid"
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
		expected Question
	}{
		"none": {
			quiz:     &Quiz{},
			expected: nil,
			index:    15,
		},
		"first": {
			quiz: &Quiz{MultipleChoiceQuestions: []*MultipleChoiceQuestion{
				{
					BaseQuestion: BaseQuestion{
						BaseObject: BaseObject{ID: uuid.MustParse("c5803f9a-a584-409a-9c06-7e66a37a959f")},
						Title:      "First question!",
						Order:      0,
					},
				},
			}},
			expected: MultipleChoiceQuestion{
				BaseQuestion: BaseQuestion{
					BaseObject: BaseObject{ID: uuid.MustParse("c5803f9a-a584-409a-9c06-7e66a37a959f")},
					Title:      "First question!",
					Order:      0,
				},
			},
			index: 0,
		},
		"second": {
			quiz: &Quiz{MultipleChoiceQuestions: []*MultipleChoiceQuestion{
				{
					BaseQuestion: BaseQuestion{
						BaseObject: BaseObject{ID: uuid.MustParse("41dbef92-9bdb-4714-8dd8-a6163ad382a7")},
						Title:      "Second question!",
						Order:      1,
					},
				},
				{
					BaseQuestion: BaseQuestion{
						BaseObject: BaseObject{ID: uuid.MustParse("8e60fff6-94c6-4e8b-8b36-6c949bcaf15d")},
						Title:      "First question!",
						Order:      0,
					},
				},
				{
					BaseQuestion: BaseQuestion{
						BaseObject: BaseObject{ID: uuid.MustParse("b582db7d-80bb-4441-a449-f1af7c1a6e94")},
						Title:      "Third question!",
						Order:      2,
					},
				},
			}},
			expected: MultipleChoiceQuestion{
				BaseQuestion: BaseQuestion{
					BaseObject: BaseObject{ID: uuid.MustParse("41dbef92-9bdb-4714-8dd8-a6163ad382a7")},
					Title:      "Second question!",
					Order:      1,
				},
			},
			index: 1,
		},
	}

	for name, testData := range tests {
		testData := testData
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			// Act
			result, ok := testData.quiz.GetQuestion(testData.index)

			// Assert
			assert.Equal(t, testData.expected != nil, ok)

			if ok {
				assert.Equal(t, testData.expected.GetBaseQuestion().Title, result.GetBaseQuestion().Title)
			}
		})
	}
}

func TestQuiz_GetNextQuestion_ReturnsExpectedValue(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		quiz            *Quiz
		currentQuestion uuid.UUID
		expected        Question
	}{
		"no current, no next": {
			quiz:            &Quiz{},
			expected:        nil,
			currentQuestion: uuid.Nil,
		},
		"no current, next": {
			quiz: &Quiz{
				MultipleChoiceQuestions: []*MultipleChoiceQuestion{
					{
						BaseQuestion: BaseQuestion{
							BaseObject: BaseObject{ID: uuid.MustParse("c5803f9a-a584-409a-9c06-7e66a37a959f")},
							Title:      "First question!",
							Order:      0,
						},
					},
				},
			},
			expected: MultipleChoiceQuestion{
				BaseQuestion: BaseQuestion{
					BaseObject: BaseObject{ID: uuid.MustParse("c5803f9a-a584-409a-9c06-7e66a37a959f")},
					Title:      "First question!",
					Order:      0,
				},
			},
			currentQuestion: uuid.Nil,
		},
		"current, no next": {
			quiz: &Quiz{
				MultipleChoiceQuestions: []*MultipleChoiceQuestion{
					{
						BaseQuestion: BaseQuestion{
							BaseObject: BaseObject{ID: uuid.MustParse("c5803f9a-a584-409a-9c06-7e66a37a959f")},
							Title:      "First question!",
							Order:      0,
						},
					},
				},
			},
			expected:        nil,
			currentQuestion: uuid.MustParse("c5803f9a-a584-409a-9c06-7e66a37a959f"),
		},
		"current and next": {
			quiz: &Quiz{
				MultipleChoiceQuestions: []*MultipleChoiceQuestion{
					{
						BaseQuestion: BaseQuestion{
							BaseObject: BaseObject{ID: uuid.MustParse("c5803f9a-a584-409a-9c06-7e66a37a959f")},
							Title:      "First question!",
							Order:      1,
						},
					},
					{
						BaseQuestion: BaseQuestion{
							BaseObject: BaseObject{ID: uuid.MustParse("dfd2294f-a2bf-45bf-9522-982b2fd056a6")},
							Title:      "Second question!",
							Order:      0,
						},
					},
				},
			},
			expected: MultipleChoiceQuestion{
				BaseQuestion: BaseQuestion{
					BaseObject: BaseObject{ID: uuid.MustParse("c5803f9a-a584-409a-9c06-7e66a37a959f")},
					Title:      "First question!",
					Order:      1,
				},
			},
			currentQuestion: uuid.MustParse("dfd2294f-a2bf-45bf-9522-982b2fd056a6"),
		},
	}

	for name, testData := range tests {
		testData := testData
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			// Act
			result, ok := testData.quiz.GetNextQuestion(testData.currentQuestion)

			// Assert
			assert.Equal(t, testData.expected != nil, ok)

			if ok {
				assert.Equal(t, testData.expected.GetBaseQuestion().Title, result.GetBaseQuestion().Title)
			}
		})
	}
}
