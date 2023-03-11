package inputs

import (
	"github.com/google/uuid"
	"github.com/survivorbat/qq.maarten.dev/server/domain"
)

type QuestionOption struct {
	TextOption string `json:"textOption" example:"Rome"` // desc: Only one value must be filled in
	Answer     bool   `json:"answer" example:"true"`     // desc: Marks this option as the answer, should only be used once
}

type MultipleChoiceQuestion struct {
	Title             string `json:"title" binding:"required,min=3,max=30" example:"What is the best city?"`
	Description       string `json:"description" example:"Subjective, but whatever ;)"`
	DurationInSeconds uint   `json:"durationInSeconds" binding:"required,min=5,max=60" example:"15"`
	Category          string `json:"category" binding:"required,min=3" example:"Geography"`
	Order             uint   `json:"order" binding:"required" example:"0"` // desc: Determines the order of this question in the quiz

	Options []*QuestionOption `json:"options" binding:"required,gte=2,lte=4"`
}

func (m MultipleChoiceQuestion) hasOneAnswer() bool {
	var foundAnswer bool

	for _, option := range m.Options {
		if option.Answer {
			// Only one answer allowed
			if foundAnswer {
				return false
			}

			foundAnswer = true
		}
	}

	return foundAnswer
}

func (m MultipleChoiceQuestion) IsValid() (bool, any, string, string, string, string) {
	if !m.hasOneAnswer() {
		return true, nil, "Options", "Options", "hasOneAnswer", "must have exactly one answer"
	}

	return false, "", "", "", "", ""
}

type Quiz struct {
	Name                    string                    `json:"name" binding:"required,min=3,max=30" example:"My awesome quiz"`
	Description             string                    `json:"description" binding:"omitempty,max=250" example:"This is going to be amazing"`
	MultipleChoiceQuestions []*MultipleChoiceQuestion `json:"multipleChoiceQuestions" binding:"lte=20"`
}

func (q Quiz) IsValid() (bool, any, string, string, string, string) {
	if !q.hasValidOrder() {
		return true, nil, "Order", "Order", "hasValidOrder", "Invalid order"
	}

	if !q.hasAnyQuestions() {
		return true, nil, "Questions", "questions", "hasAnyQuestions", "No questions"
	}

	return false, "", "", "", "", ""
}

func (q Quiz) hasAnyQuestions() bool {
	return len(q.MultipleChoiceQuestions) > 0
}

// hasValidOrder verifies whether the questions are ordered correctly
func (q Quiz) hasValidOrder() bool {
	var count uint
	var control uint

	for index, question := range q.MultipleChoiceQuestions {
		count += question.Order
		control += uint(index)
	}

	return count == control
}

// NewUuid may be overwritten in tests
var NewUuid = uuid.New

// ToDomain is tested through the routes test

func (q Quiz) ToDomain() *domain.Quiz {
	mcQuestions := make([]*domain.MultipleChoiceQuestion, len(q.MultipleChoiceQuestions))
	for index, mcQuestion := range q.MultipleChoiceQuestions {
		result := &domain.MultipleChoiceQuestion{
			BaseQuestion: domain.BaseQuestion{
				Title:             mcQuestion.Title,
				Description:       mcQuestion.Description,
				DurationInSeconds: mcQuestion.DurationInSeconds,
				Category:          mcQuestion.Category,
				Order:             mcQuestion.Order,
			},
			Options: make([]*domain.QuestionOption, len(mcQuestion.Options)),
		}

		for index, mcOption := range mcQuestion.Options {
			result.Options[index] = &domain.QuestionOption{
				BaseObject: domain.BaseObject{ID: NewUuid()},
				TextOption: mcOption.TextOption,
			}

			if mcOption.Answer {
				result.AnswerID = result.Options[index].ID
			}
		}

		mcQuestions[index] = result
	}

	return &domain.Quiz{
		Name:                    q.Name,
		Description:             q.Description,
		MultipleChoiceQuestions: mcQuestions,
	}
}
