package domain

import (
	"github.com/google/uuid"
	"time"
)

type Question interface {
	GetBaseQuestion() BaseQuestion
	GetType() QuestionType
}

type BaseQuestion struct {
	Title       string
	Description string
	Time        time.Duration
	Category    string

	QuizID uuid.UUID `json:"quizID"`
	Quiz   *Quiz     `json:"quiz" gorm:"foreignKey:QuizID"`
}
