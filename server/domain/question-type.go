package domain

const QuestionTypeMC QuestionType = "multiple choice"

type QuestionType string

func (q QuestionType) IsValid() bool {
	switch q {
	case QuestionTypeMC:
		return true
	}

	return false
}
