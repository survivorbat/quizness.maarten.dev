package outputs

import (
	"github.com/google/uuid"
	"github.com/survivorbat/qq.maarten.dev/server/domain"
)

func NewPublicQuiz(quiz *domain.Quiz) *OutputQuiz {
	result := &OutputQuiz{
		ID:                      quiz.ID,
		Name:                    quiz.Name,
		Description:             quiz.Description,
		MultipleChoiceQuestions: make([]*OutputMultipleChoiceQuestion, len(quiz.MultipleChoiceQuestions)),
	}

	for index, mc := range quiz.MultipleChoiceQuestions {
		result.MultipleChoiceQuestions[index] = &OutputMultipleChoiceQuestion{
			ID:                mc.ID,
			Title:             mc.Title,
			Description:       mc.Description,
			DurationInSeconds: mc.DurationInSeconds,
			Category:          mc.Category,
			Order:             mc.Order,
			Options:           mc.Options,
		}
	}

	return result
}

type OutputQuiz struct {
	ID uuid.UUID `json:"id"`

	Name        string `json:"name" example:"Daniel's funky quiz'"`     // desc: Can be anything
	Description string `json:"description" example:"My first attempt!"` // desc: Ditto

	MultipleChoiceQuestions []*OutputMultipleChoiceQuestion `json:"multipleChoiceQuestions,omitempty"`
}
