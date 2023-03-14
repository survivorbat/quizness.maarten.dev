package domain

import "github.com/google/uuid"

type GameAnswers []*GameAnswer

func (g GameAnswers) Contains(questionID uuid.UUID, playerID uuid.UUID) bool {
	for _, answer := range g {
		if answer.QuestionID == questionID && answer.PlayerID == playerID {
			return true
		}
	}

	return false
}

type GameAnswer struct {
	BaseObject

	PlayerID uuid.UUID `json:"playerID"  example:"00000000-0000-0000-0000-000000000000"`
	Player   *Player   `json:"player" gorm:"foreignKey:PlayerID"`

	GameID uuid.UUID `json:"gameID"  example:"00000000-0000-0000-0000-000000000000"`
	Game   *Game     `json:"game" gorm:"foreignKey:GameID"`

	QuestionID uuid.UUID `json:"questionID"  example:"00000000-0000-0000-0000-000000000000"`
	OptionID   uuid.UUID `json:"optionID"  example:"00000000-0000-0000-0000-000000000000"`
}
