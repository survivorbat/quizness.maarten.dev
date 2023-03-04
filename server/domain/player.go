package domain

import "github.com/google/uuid"

type Player struct {
	BaseObject

	NickName string `json:"nickName"`

	GameID uuid.UUID `json:"gameID"`
	Game   *Game     `json:"game" gorm:"GameID"`
}
