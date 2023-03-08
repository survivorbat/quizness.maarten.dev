package domain

import "github.com/google/uuid"

// QuestionOption represents a potential solution for a quiz question, can be anything.
type QuestionOption struct {
	BaseObject
	MultipleChoiceQuestion   MultipleChoiceQuestion `json:"-" gorm:"foreignKey:MultipleChoiceQuestionID"`
	MultipleChoiceQuestionID *uuid.UUID             `json:"-"`

	TextOption string `json:"textOption" example:"Haarlem"` // desc: A textual option for this question

	// Expand for more option types
	// ImageOption
	// ...
}
