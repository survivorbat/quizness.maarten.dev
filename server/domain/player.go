package domain

import (
	"fmt"
	"github.com/google/uuid"
)

type Players []*Player

func (p Players) Contains(playerID uuid.UUID) bool {
	for _, player := range p {
		if player.ID == playerID {
			return true
		}
	}

	return false
}

// Player is a simple anonymous user that joins a game, also gets a random nickname assigned
type Player struct {
	BaseObject

	Nickname string `json:"nickname" example:"Adorable Beaver"` // desc: Randomly assigned nickname, to avoid naughty words

	GameID uuid.UUID `json:"gameID" example:"00000000-0000-0000-0000-000000000000"` // desc: The game this player belongs to
	Game   *Game     `json:"-" gorm:"foreignKey:GameID"`
}

// GenerateNickname overwrites the creator's nickname using a random prefix and suffix
func (c *Player) GenerateNickname() {
	prefix := namePrefixes[randomGenerator.Intn(len(namePrefixes))]
	suffix := nameSuffixes[randomGenerator.Intn(len(nameSuffixes))]
	c.Nickname = fmt.Sprintf("%s %s", prefix, suffix)
}
