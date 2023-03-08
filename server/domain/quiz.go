package domain

import "github.com/google/uuid"

type Quiz struct {
	BaseObject

	Name        string `json:"name"`
	Description string `json:"description"`

	CreatorID uuid.UUID `json:"creatorID"`
	Creator   *Creator  `json:"creator" gorm:"foreignKey:CreatorID"`

	MultipleChoiceQuestions []*MultipleChoiceQuestion `json:"multipleChoiceQuestions,omitempty" gorm:"foreignKey:QuizID;constraint:OnDelete:CASCADE"`
	Games                   []*Game                   `json:"games" gorm:"foreignKey:QuizID;constraint:OnDelete:CASCADE"`
}
