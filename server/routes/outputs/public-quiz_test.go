package outputs

import (
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
		expected := &OutputMultipleChoiceQuestion{
			ID:                quiz.MultipleChoiceQuestions[0].ID,
			Title:             quiz.MultipleChoiceQuestions[0].Title,
			Description:       quiz.MultipleChoiceQuestions[0].Description,
			DurationInSeconds: quiz.MultipleChoiceQuestions[0].DurationInSeconds,
			Category:          quiz.MultipleChoiceQuestions[0].Category,
			Order:             quiz.MultipleChoiceQuestions[0].Order,
			Options:           quiz.MultipleChoiceQuestions[0].Options,
		}
		assert.Equal(t, expected, result.MultipleChoiceQuestions[0])
	}
}
