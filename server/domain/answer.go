package domain

import "github.com/google/uuid"

type GameAnswer struct {
	BaseObject

	PlayerID uuid.UUID `json:"playerID"  example:"00000000-0000-0000-0000-000000000000"`
	Player   *Player   `json:"player" gorm:"foreignKey:PlayerID"`

	GameID uuid.UUID `json:"gameID"  example:"00000000-0000-0000-0000-000000000000"`
	Game   *Game     `json:"game" gorm:"foreignKey:GameID"`

	QuestionID uuid.UUID `json:"questionID"  example:"00000000-0000-0000-0000-000000000000"`
	OptionID   uuid.UUID `json:"optionID"  example:"00000000-0000-0000-0000-000000000000"`
}
