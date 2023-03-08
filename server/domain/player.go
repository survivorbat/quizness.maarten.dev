package domain

import (
	"fmt"
	"github.com/google/uuid"
)

type Player struct {
	BaseObject

	Nickname string `json:"nickname"`

	GameID uuid.UUID `json:"gameID"`
	Game   *Game     `json:"game" gorm:"foreignKey:GameID"`
}

func (c *Player) GenerateNickname() {
	prefix := namePrefixes[randomGenerator.Intn(len(namePrefixes))]
	suffix := nameSuffixes[randomGenerator.Intn(len(nameSuffixes))]
	c.Nickname = fmt.Sprintf("%s %s", prefix, suffix)
}
