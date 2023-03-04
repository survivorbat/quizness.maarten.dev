package domain

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type Game struct {
	BaseObject

	QuizID uuid.UUID `json:"quizID"`
	Quiz   *Quiz     `json:"quiz" gorm:"QuizID"`

	Code        string    `json:"code"`
	PlayerLimit uint      `json:"limit"`
	Players     []*Player `json:"players,omitempty" gorm:"foreignKey:GameID"`
	StartTime   time.Time `json:"startTime"`
	FinishTime  time.Time `json:"finishTime"`
}

// playerCompare is used in the containsWithKey function
func playerCompare(p *Player) string {
	return p.NickName
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

// Start starts the game
func (q *Game) Start() {
	q.StartTime = time.Now()
}
