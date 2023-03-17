package coordinator

type PlayerAction string

const AnswerAction PlayerAction = "answer"

func (c PlayerAction) IsValid() bool {
	switch c {
	case AnswerAction:
		return true
	default:
		return false
	}
}

type CreatorAction string

const NextQuestionAction CreatorAction = "next"
const FinishGameAction CreatorAction = "finish"

func (c CreatorAction) IsValid() bool {
	switch c {
	case FinishGameAction, NextQuestionAction:
		return true
	default:
		return false
	}
}
