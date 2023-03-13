package domain

import (
	"errors"
	"github.com/google/uuid"
	"math/rand"
	"strings"
	"time"
)

// Game is an occurrence of a quiz, a quiz can be conducted multiple times
type Game struct {
	BaseObject

	QuizID uuid.UUID `json:"quizID" example:"00000000-0000-0000-0000-000000000000"`
	Quiz   *Quiz     `json:"-" gorm:"foreignKey:QuizID"`

	Code        string `json:"code" example:"KO384B"` // desc: The 'join' code for new players
	PlayerLimit uint   `json:"playerLimit"`           // desc: The max amount of players that may join this game

	CurrentQuestion uuid.UUID `json:"currentQuestion" example:"00000000-0000-0000-0000-000000000000"` // desc: The current question

	Players []*Player     `json:"players" gorm:"foreignKey:GameID;constraint:OnDelete:CASCADE"`
	Answers []*GameAnswer `json:"answers" gorm:"foreignKey:GameID;constraint:OnDelete:CASCADE"`

	StartTime  time.Time `json:"startTime"`  // desc: The time that this game started
	FinishTime time.Time `json:"finishTime"` // desc: The time that this game ended
}

// IsInProgress returns whether the game has started
func (g *Game) IsInProgress() bool {
	return !g.StartTime.IsZero() && g.FinishTime.IsZero()
}

// playerCompare is used in the containsWithKey function
func playerCompare(p *Player) string {
	return p.Nickname
}

// PlayerJoin adds a player to the game, returns an error if a player with the nickname
// is already present
func (g *Game) PlayerJoin(player *Player) error {
	if _, ok := containsWithKey(player, g.Players, playerCompare); ok {
		return errors.New("player is already in this game")
	}

	g.Players = append(g.Players, player)
	return nil
}

// PlayerLeave removes a player from the game, returns an error if the player is not
// present
func (g *Game) PlayerLeave(player *Player) error {
	index, ok := containsWithKey(player, g.Players, playerCompare)
	if !ok {
		return errors.New("player is not in this game")
	}

	g.Players = append(g.Players[:index], g.Players[index+1:]...)
	return nil
}

// Start starts the game and sets the code
func (g *Game) Start() error {
	if !g.StartTime.IsZero() {
		return errors.New("game has already started")
	}

	g.StartTime = time.Now()

	code := make([]string, 6)
	for i := range code {
		code[i] = codeChars[rand.Intn(len(codeChars))]
	}
	g.Code = strings.Join(code, "")

	return nil
}

// Finish ends the game
func (g *Game) Finish() error {
	if g.StartTime.IsZero() {
		return errors.New("game has not started")
	}

	if !g.FinishTime.IsZero() {
		return errors.New("game has already finished")
	}

	g.FinishTime = time.Now()
	return nil
}
