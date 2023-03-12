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

	Players []*Player `json:"-" gorm:"foreignKey:GameID;constraint:OnDelete:CASCADE"`

	StartTime  time.Time `json:"startTime"`  // desc: The time that this game started
	FinishTime time.Time `json:"finishTime"` // desc: The time that this game ended
}

// playerCompare is used in the containsWithKey function
func playerCompare(p *Player) string {
	return p.Nickname
}

// PlayerJoin adds a player to the game, returns an error if a player with the nickname
// is already present
func (q *Game) PlayerJoin(player *Player) error {
	if _, ok := containsWithKey(player, q.Players, playerCompare); ok {
		return errors.New("player is already in this game")
	}

	q.Players = append(q.Players, player)
	return nil
}

// PlayerLeave removes a player from the game, returns an error if the player is not
// present
func (q *Game) PlayerLeave(player *Player) error {
	index, ok := containsWithKey(player, q.Players, playerCompare)
	if !ok {
		return errors.New("player is not in this game")
	}

	q.Players = append(q.Players[:index], q.Players[index+1:]...)
	return nil
}

// Start starts the game and sets the code
func (q *Game) Start() error {
	if !q.StartTime.IsZero() {
		return errors.New("game has already started")
	}

	q.StartTime = time.Now()

	code := make([]string, 6)
	for i := range code {
		code[i] = codeChars[rand.Intn(len(codeChars))]
	}
	q.Code = strings.Join(code, "")

	return nil
}

// Finish ends the game
func (q *Game) Finish() error {
	if q.StartTime.IsZero() {
		return errors.New("game has not started")
	}

	if !q.FinishTime.IsZero() {
		return errors.New("game has already finished")
	}

	q.FinishTime = time.Now()
	return nil
}
