package outputs

import (
	"github.com/google/uuid"
	"github.com/survivorbat/qq.maarten.dev/server/domain"
)

type OutputMultipleChoiceQuestion struct {
	ID                uuid.UUID                `json:"id" example:"00000000-0000-0000-0000-000000000000"`
	Title             string                   `json:"title" example:"What is 5+5?"`
	Description       string                   `json:"description" example:"We want to test your math skills for no apparent reason"`
	DurationInSeconds uint                     `json:"durationInSeconds" example:"30"`
	Category          string                   `json:"category" example:"Geography"`
	Order             uint                     `json:"order" example:"2"`
	Options           []*domain.QuestionOption `json:"options"`
}
