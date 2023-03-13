package domain

import (
	"github.com/google/uuid"
)

// Question is used to pass around questions without a specific type
type Question interface {
	GetBaseQuestion() BaseQuestion
}

// BaseQuestion contains fields that every question should contain, allows us to embed it in other questions
type BaseQuestion struct {
	BaseObject
	Title             string `json:"title" example:"What is 5+5?"`                                                  // desc: Can be anything
	Description       string `json:"description" example:"We want to test your math skills for no apparent reason"` // desc: May be empty
	DurationInSeconds uint   `json:"durationInSeconds" example:"30"`                                                // desc: Players have this amount of seconds to
	Category          string `json:"category" example:"Geography"`                                                  // desc: Allows you to group questions together, does not affect order
	Order             uint   `json:"order" example:"2"`                                                             // desc: The order of this question in the quiz

	QuizID uuid.UUID `json:"quizID" example:"00000000-0000-0000-0000-000000000000"` // desc: The quiz this question belongs to
	Quiz   *Quiz     `json:"-" gorm:"foreignKey:QuizID"`
}

func (b BaseQuestion) GetBaseQuestion() BaseQuestion {
	return b
}
