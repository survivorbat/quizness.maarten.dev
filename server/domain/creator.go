package domain

import (
	"fmt"
	"github.com/AvraamMavridis/randomcolor"
)

// Creator is a user that makes quiz and conducts games
type Creator struct {
	BaseObject

	Color           string `json:"color" example:"#220022"`                          // desc: Randomly assigned color
	BackgroundColor string `json:"backgroundColor" example:"#220022"`                // desc: Randomly assigned color
	Nickname        string `json:"nickname" gorm:"unique" example:"Adorable Beaver"` // desc: Randomly assigned nickname, to avoid naughty words

	// Never expose this
	AuthID string `json:"-" gorm:"unique"`

	Quizzes []*Quiz `json:"-" gorm:"foreignKey:CreatorID;constraint:OnDelete:CASCADE"`
}

// GenerateNickname overwrites the creator's nickname using a random prefix and suffix
func (c *Creator) GenerateNickname() {
	prefix := namePrefixes[randomGenerator.Intn(len(namePrefixes))]
	suffix := nameSuffixes[randomGenerator.Intn(len(nameSuffixes))]
	c.Nickname = fmt.Sprintf("%s %s", prefix, suffix)
}

// GenerateColors overwrites the creator's colors
func (c *Creator) GenerateColors() {
	c.Color = randomcolor.GetRandomColorInHex()
	c.BackgroundColor = randomcolor.GetRandomColorInHex()
}
