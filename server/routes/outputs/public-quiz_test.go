package outputs

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/survivorbat/qq.maarten.dev/server/domain"
	"testing"
)

func TestNewPublicQuiz_ReturnsExpectedOutput(t *testing.T) {
	t.Parallel()
	// Arrange
	quiz := &domain.Quiz{
		BaseObject:  domain.BaseObject{ID: uuid.MustParse("00000000-0000-0000-0000-000000000000")},
		Name:        "abc",
		Description: "def",
		MultipleChoiceQuestions: []*domain.MultipleChoiceQuestion{
			{BaseQuestion: domain.BaseQuestion{Title: "ghi"}},
		},
	}

	// Act
	result := NewPublicQuiz(quiz)

	// Assert
	assert.Equal(t, quiz.ID, result.ID)
	assert.Equal(t, quiz.Name, result.Name)
	assert.Equal(t, quiz.Description, result.Description)

	if assert.Len(t, quiz.MultipleChoiceQuestions, 1) {
		expected, _ := NewPublicQuestion(quiz.MultipleChoiceQuestions[0])
		assert.Equal(t, json.RawMessage(expected), result.MultipleChoiceQuestions[0])
	}
}
