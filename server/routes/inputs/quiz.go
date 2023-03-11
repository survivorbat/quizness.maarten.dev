package inputs

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
