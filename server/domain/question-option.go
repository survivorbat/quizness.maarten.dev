package domain

import "github.com/google/uuid"

// QuestionOption represents a potential solution for a quiz question, can be anything.
type QuestionOption struct {
	BaseObject
	MultipleChoiceQuestionID *uuid.UUID `json:"-"`

	TextOption string `json:"textOption"`

	// Expand for more option types
	// ImageOption
	// ...
}
