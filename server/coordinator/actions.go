package coordinator

type PlayerAction string

const JoinAction PlayerAction = "join"
const AnswerAction PlayerAction = "answer"
const LeaveAction PlayerAction = "leave"

func (c PlayerAction) IsValid() bool {
	switch c {
	case JoinAction, AnswerAction, LeaveAction:
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
