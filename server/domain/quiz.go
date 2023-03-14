package domain

import (
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

// GetQuestion retrieves a question based on the Order of a question in the list, will return
// uuid.Nil if not found
func (q *Quiz) GetQuestion(order uint) (Question, bool) {
	for _, question := range q.MultipleChoiceQuestions {
		if question.Order == order {
			return question, true
		}
	}

	return nil, false
}

func (q *Quiz) GetNextQuestion(current uuid.UUID) (Question, bool) {
	// If no current is found, take the first one
	if current == uuid.Nil {
		return q.GetQuestion(0)
	}

	var currentQuestion uint
	for _, question := range q.MultipleChoiceQuestions {
		if question.ID == current {
			currentQuestion = question.Order
			break
		}
	}

	return q.GetQuestion(currentQuestion + 1)
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
