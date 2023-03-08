package domain

import (
	"fmt"
)

// Creator is a user that makes quiz and conducts games
type Creator struct {
	BaseObject

	Nickname string `json:"nickname" gorm:"unique" example:"Adorable Beaver"` // desc: Randomly assigned nickname, to avoid naughty words

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
