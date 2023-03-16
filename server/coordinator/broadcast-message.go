package coordinator

import "github.com/google/uuid"

type BroadcastType string

const NextQuestionType BroadcastType = "next"
const FinishGameType BroadcastType = "finish"
const PlayerAnsweredType BroadcastType = "answered"

type BroadcastMessage struct {
	Type       BroadcastType `json:"type"`
	QuestionID uuid.UUID     `json:"questionID"`
	PlayerID   uuid.UUID     `json:"playerID"`
}
