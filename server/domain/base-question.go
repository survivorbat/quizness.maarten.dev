package domain

import (
	"github.com/google/uuid"
)

// BaseQuestion contains fields that every question should contain
type BaseQuestion struct {
	BaseObject
	Title             string `json:"title"`
	Description       string `json:"description"`
	DurationInSeconds uint   `json:"durationInSeconds"`
	Category          string `json:"category"`
	Order             int    `json:"order"`

	QuizID uuid.UUID `json:"quizID"`
	Quiz   *Quiz     `json:"quiz" gorm:"foreignKey:QuizID"`
}
