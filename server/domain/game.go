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
	CurrentDeadline time.Time `json:"currentDeadline"`                                                // desc: Past this deadline, no answers may be submitted

	Players Players     `json:"players" gorm:"foreignKey:GameID;constraint:OnDelete:CASCADE"`
	Answers GameAnswers `json:"answers" gorm:"foreignKey:GameID;constraint:OnDelete:CASCADE"`

	StartTime  time.Time `json:"startTime"`  // desc: The time that this game started
	FinishTime time.Time `json:"finishTime"` // desc: The time that this game ended
}

func (g *Game) GetCurrentQuestion() (Question, bool) {
	// Give up quick
	if g.CurrentQuestion == uuid.Nil {
		return nil, false
	}

	for _, question := range g.Quiz.MultipleChoiceQuestions {
		if question.ID == g.CurrentQuestion {
			return question, true
		}
	}

	return nil, false
}

func (g *Game) IsInProgress() bool {
	return !g.StartTime.IsZero() && g.FinishTime.IsZero()
}

func (g *Game) IsOpenForPlayers() bool {
	return !g.StartTime.IsZero() && g.CurrentQuestion == uuid.Nil && g.FinishTime.IsZero()
}

// Start starts the game and sets the code
func (g *Game) Start() error {
	if !g.StartTime.IsZero() {
		return errors.New("game has already started")
	}

	if g.Quiz.CountQuestions() == 0 {
		return errors.New("no questions defined")
	}

	g.StartTime = time.Now()

	code := make([]string, 6)
	for i := range code {
		code[i] = codeChars[rand.Intn(len(codeChars))]
	}
	g.Code = strings.Join(code, "")

	return nil
}

func (g *Game) Next() error {
	if !g.IsInProgress() {
		return errors.New("game is not in progress")
	}

	if len(g.Players) < 2 {
		return errors.New("can only start with 2 or more players")
	}

	if !g.CurrentDeadline.IsZero() && time.Now().Before(g.CurrentDeadline) {
		return errors.New("deadline has not passed")
	}

	nextQuestion, ok := g.Quiz.GetNextQuestion(g.CurrentQuestion)

	if !ok {
		return errors.New("no more questions")
	}

	g.CurrentQuestion = nextQuestion.GetBaseQuestion().ID
	g.CurrentDeadline = time.Now().Add(time.Duration(nextQuestion.GetBaseQuestion().DurationInSeconds) * time.Second)

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

	if !g.CurrentDeadline.IsZero() && time.Now().Before(g.CurrentDeadline) {
		return errors.New("deadline has not passed")
	}

	g.FinishTime = time.Now()
	return nil
}

func (g *Game) AnswerQuestion(player uuid.UUID, question uuid.UUID, optionID uuid.UUID) (*GameAnswer, error) {
	if g.CurrentQuestion != question {
		return nil, errors.New("not the current question")
	}

	if time.Now().After(g.CurrentDeadline) {
		return nil, errors.New("deadline passed")
	}

	if !g.Players.Contains(player) {
		return nil, errors.New("player not in game")
	}

	if g.Answers.Contains(question, player) {
		return nil, errors.New("player has already submitted an answer")
	}

	answer := &GameAnswer{
		PlayerID:   player,
		GameID:     g.ID,
		QuestionID: question,
		OptionID:   optionID,
	}

	g.Answers = append(g.Answers, answer)

	return answer, nil
}
