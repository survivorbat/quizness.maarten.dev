package domain

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	namePrefixes = []string{
		"Blazing",
		"Crazy",
		"Amazing",
		"Furry",
		"Hairy",
		"Gentle",
		"Strong",
		"Adorable",
		"Adventurous",
		"Dizzy",
		"Cute",
		"Clever",
		"Colorful",
		"Busy",
		"Brave",
		"Calm",
		"Brainy",
		"Bright",
		"Concerned",
		"Curious",
	}
	nameSuffixes = []string{
		"Beaver",
		"Olive",
		"Apple",
		"Pear",
		"Potato",
		"Cat",
		"Puppy",
		"Whiteboard",
		"Bear",
		"Fox",
		"Elephant",
		"Fox",
		"Raccoon",
		"Gazelle",
		"Phone",
		"Deer",
		"Lemonade",
		"Joker",
		"Rabbit",
		"Zombie",
		"Skeleton",
	}
	randSource      = rand.NewSource(time.Now().Unix())
	randomGenerator = rand.New(randSource)
)

type Creator struct {
	BaseObject

	Nickname string `json:"nickname" gorm:"unique"`

	// Never expose this
	AuthID string `json:"-" gorm:"unique"`

	Quizzes []*Quiz `json:"quizzes,omitempty" gorm:"foreignKey:CreatorID"`
}

func (c *Creator) GenerateNickname() {
	prefix := namePrefixes[randomGenerator.Intn(len(namePrefixes))]
	suffix := nameSuffixes[randomGenerator.Intn(len(nameSuffixes))]
	c.Nickname = fmt.Sprintf("%s %s", prefix, suffix)
}
