package coordinator

import (
	"github.com/google/uuid"
)

type BroadcastType string

// NextQuestionType is used to progress the game
const NextQuestionType BroadcastType = "next"

// FinishGameType is used to end a current game
const FinishGameType BroadcastType = "finish"

// PlayerAnsweredType is used to indicate another player answered a question
const PlayerAnsweredType BroadcastType = "answered"

// ParticipantsType is used to broadcast the current participants and the creator
const ParticipantsType BroadcastType = "participants"

type BroadcastMessage struct {
	Type BroadcastType `json:"type"`

	// NextQuestionType type
	NextQuestionContent *nextQuestionContent `json:"nextQuestionContent,omitempty"`

	// PlayerAnsweredType
	PlayerAnsweredContent *playerAnsweredContent `json:"playerAnsweredContent,omitempty"`

	// ParticipantsType type
	ParticipantsContent *participantsContent `json:"participantsContent,omitempty"`
}

type participantsContent struct {
	Creator *participant   `json:"creator"`
	Players []*participant `json:"players"`
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
