package outputs

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/survivorbat/qq.maarten.dev/server/domain"
	"testing"
)

func TestNewPublicQuestion_ReturnsExpectedValue(t *testing.T) {
	t.Parallel()
	tests := map[domain.QuestionType]struct {
		input    *domain.MultipleChoiceQuestion
		expected OutputMultipleChoiceQuestion
	}{
		domain.TypeMultipleChoice: {
			input: &domain.MultipleChoiceQuestion{
				BaseQuestion: domain.BaseQuestion{
					BaseObject:        domain.BaseObject{ID: uuid.MustParse("751f3398-e129-463d-b76d-bae36c0321c3")},
					Title:             "abc",
					Description:       "def",
					DurationInSeconds: 15,
					Category:          "ghi",
					Order:             3,
				},
				Options: []*domain.QuestionOption{
					{TextOption: "1"},
					{TextOption: "2"},
				},
			},
			expected: OutputMultipleChoiceQuestion{
				ID:                uuid.MustParse("751f3398-e129-463d-b76d-bae36c0321c3"),
				Title:             "abc",
				Description:       "def",
				DurationInSeconds: 15,
				Category:          "ghi",
				Order:             3,
				Options: []*domain.QuestionOption{
					{TextOption: "1"},
					{TextOption: "2"},
				},
			},
		},
	}

	for name, testData := range tests {
		testData := testData
		t.Run(string(name), func(t *testing.T) {
			t.Parallel()
			// Act
			result, err := NewPublicQuestion(testData.input)

			// Assert
			assert.NoError(t, err)

			expected, _ := json.Marshal(testData.expected)
			assert.Equal(t, string(expected), string(result))
		})
	}
}

type noQuestion struct {
	domain.BaseQuestion
}

func (n noQuestion) GetType() domain.QuestionType {
	return ""
}

func TestNewPublicQuestion_ReturnsErrorOnUnknownType(t *testing.T) {
	t.Parallel()
	// Arrange
	question := &noQuestion{}

	// Act
	result, err := NewPublicQuestion(question)

	// Assert
	assert.Empty(t, result)
	assert.ErrorContains(t, err, "not found")
}
