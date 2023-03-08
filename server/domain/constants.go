package domain

import (
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
