package domain

import "github.com/google/uuid"

// MultipleChoiceQuestion is one of the potential questions a user can add to their quiz
type MultipleChoiceQuestion struct {
	BaseQuestion

	AnswerID uuid.UUID `json:"answerID" example:"00000000-0000-0000-0000-000000000000"`

	Options []*QuestionOption `json:"options" gorm:"foreignKey:MultipleChoiceQuestionID;constraint:OnDelete:CASCADE"`
}

func (m MultipleChoiceQuestion) GetType() QuestionType {
	return TypeMultipleChoice
}
