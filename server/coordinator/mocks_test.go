package coordinator

import (
	"github.com/google/uuid"
	"github.com/survivorbat/qq.maarten.dev/server/domain"
	"github.com/survivorbat/qq.maarten.dev/server/services"
)

type MockGameService struct {
	services.GameService

	answerQuestionCalledWithGame     *domain.Game
	answerQuestionCalledWithQuestion uuid.UUID
	answerQuestionCalledWithPlayer   uuid.UUID
	answerQuestionCalledWithOption   uuid.UUID
	answerQuestionReturns            error

	getByIDReturns      *domain.Game
	getByIDReturnsError error

	nextCalledWith          *domain.Game
	nextSetsCurrentQuestion uuid.UUID
	nextReturns             error
}

func (m *MockGameService) GetByID(uuid.UUID) (*domain.Game, error) {
	return m.getByIDReturns, m.getByIDReturnsError
}

func (m *MockGameService) Next(game *domain.Game) error {
	m.nextCalledWith = game
	game.CurrentQuestion = m.nextSetsCurrentQuestion
	return m.nextReturns
}

func (m *MockGameService) AnswerQuestion(game *domain.Game, questionID uuid.UUID, playerID uuid.UUID, optionID uuid.UUID) error {
	m.answerQuestionCalledWithGame = game
	m.answerQuestionCalledWithQuestion = questionID
	m.answerQuestionCalledWithPlayer = playerID
	m.answerQuestionCalledWithOption = optionID
	return m.answerQuestionReturns
}
