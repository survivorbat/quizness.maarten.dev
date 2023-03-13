package domain

import (
	"errors"
	"github.com/google/uuid"
)

// Quiz is created by a creator and can be played multiple times in the Game object
type Quiz struct {
	BaseObject

	Name        string `json:"name" example:"Daniel's funky quiz'"`     // desc: Can be anything
	Description string `json:"description" example:"My first attempt!"` // desc: Ditto

	CreatorID uuid.UUID `json:"creatorID" example:"00000000-0000-0000-0000-000000000000"`
	Creator   *Creator  `json:"-" gorm:"foreignKey:CreatorID"`

	MultipleChoiceQuestions []*MultipleChoiceQuestion `json:"multipleChoiceQuestions,omitempty" gorm:"foreignKey:QuizID;constraint:OnDelete:CASCADE"`

	Games []*Game `json:"-" gorm:"foreignKey:QuizID;constraint:OnDelete:CASCADE"`
}

// GetQuestion retrieves a question based on the Order of a question in the list
func (q *Quiz) GetQuestion(index uint) (Question, error) {
	for _, question := range q.MultipleChoiceQuestions {
		if question.Order == index {
			return question, nil
		}
	}

	return nil, errors.New("question not found")
}

// HasGameInProgress can be used to verify whether a new game can be started
func (q *Quiz) HasGameInProgress() bool {
	for _, game := range q.Games {
		if game.IsInProgress() {
			return true
		}
	}

	return false
}
