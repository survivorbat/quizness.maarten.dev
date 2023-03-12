package domain

import (
	"math/rand"
	"strings"
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
		"Deer",
		"Lemonade",
		"Joker",
		"Rabbit",
		"Zombie",
		"Skeleton",
		"Apple",
		"Pear",
		"Banana",
		"Mango",
		"Lemon",
		"Typewriter",
	}
	codeChars       = strings.Split("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", "")
	randSource      = rand.NewSource(time.Now().Unix())
	randomGenerator = rand.New(randSource)
)
