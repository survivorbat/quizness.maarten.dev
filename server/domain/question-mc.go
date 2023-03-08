package domain

import "github.com/google/uuid"

// MultipleChoiceQuestion is one of the potential questions a user can add to their quiz
type MultipleChoiceQuestion struct {
	BaseQuestion

	// Not exposed for obvious reasons
	AnswerID uuid.UUID       `json:"-"`
	Answer   *QuestionOption `json:"-" gorm:"foreignKey:AnswerID;constraint:OnDelete:CASCADE"`

	Options []*QuestionOption `json:"options" gorm:"foreignKey:MultipleChoiceQuestionID;constraint:OnDelete:CASCADE"`
}
