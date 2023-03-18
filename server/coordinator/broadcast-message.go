package coordinator

import (
	"github.com/google/uuid"
	"time"
)

type BroadcastType string

// NextQuestionType is used to progress the game
const NextQuestionType BroadcastType = "next"

// FinishGameType is used to end a current game
const FinishGameType BroadcastType = "finish"

// PlayerAnsweredType is used to indicate another player answered a question
const PlayerAnsweredType BroadcastType = "answered"

// StateType is used to broadcast the current participants and the creator
const StateType BroadcastType = "state"

type BroadcastMessage struct {
	Type BroadcastType `json:"type"`

	// NextQuestionType type
	NextQuestionContent *nextQuestionContent `json:"nextQuestionContent,omitempty"`

	// PlayerAnsweredType
	PlayerAnsweredContent *playerAnsweredContent `json:"playerAnsweredContent,omitempty"`

	// StateType type
	StateContent *stateContent `json:"stateContent,omitempty"`
}

type stateContent struct {
	Creator         *participant   `json:"creator"`
	Players         []*participant `json:"players"`
	CurrentQuestion uuid.UUID      `json:"currentQuestion"`
	CurrentDeadline time.Time      `json:"currentDeadline"`
}

type participant struct {
	ID       uuid.UUID `json:"id"`
	Nickname string    `json:"nickname"`
}

type playerAnsweredContent struct {
	PlayerID uuid.UUID `json:"playerID"`
}

type nextQuestionContent struct {
	QuestionID uuid.UUID `json:"questionID"`
}
