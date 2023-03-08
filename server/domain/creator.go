package domain

import (
	"fmt"
)

type Creator struct {
	BaseObject

	Nickname string `json:"nickname" gorm:"unique"`

	// Never expose this
	AuthID string `json:"-" gorm:"unique"`

	Quizzes []*Quiz `json:"quizzes,omitempty" gorm:"foreignKey:CreatorID;constraint:OnDelete:CASCADE"`
}

func (c *Creator) GenerateNickname() {
	prefix := namePrefixes[randomGenerator.Intn(len(namePrefixes))]
	suffix := nameSuffixes[randomGenerator.Intn(len(nameSuffixes))]
	c.Nickname = fmt.Sprintf("%s %s", prefix, suffix)
}
