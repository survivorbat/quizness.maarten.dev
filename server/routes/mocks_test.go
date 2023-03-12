package routes

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/survivorbat/qq.maarten.dev/server/domain"
	"github.com/survivorbat/qq.maarten.dev/server/services"
)

type MockCreatorService struct {
	services.CreatorService

	getByIDCalledWith   uuid.UUID
	getByIDReturns      *domain.Creator
	getByIDReturnsError error
}

func (m *MockCreatorService) GetByID(id uuid.UUID) (*domain.Creator, error) {
	m.getByIDCalledWith = id
	return m.getByIDReturns, m.getByIDReturnsError
}

type MockGameService struct {
	services.GameService

	getByQuizCalledWith   uuid.UUID
	getByQuizReturns      []*domain.Game
	getByQuizReturnsError error

	createCalledWith *domain.Game
	createReturns    error
}

func (m *MockGameService) GetByQuiz(id uuid.UUID) ([]*domain.Game, error) {
	m.getByQuizCalledWith = id
	return m.getByQuizReturns, m.getByQuizReturnsError
}

func (m *MockGameService) Create(game *domain.Game) error {
	m.createCalledWith = game
	return m.createReturns
}

type MockQuizService struct {
	services.QuizService

	getByCreatorCalledWith   uuid.UUID
	getByCreatorReturns      []*domain.Quiz
	getByCreatorReturnsError error

	createOrUpdateCalledWith *domain.Quiz
	createOrUpdateReturns    error

	getByIdReturns      *domain.Quiz
	getByIdReturnsError error

	deleteCalledWith uuid.UUID
	deleteReturns    error
}

func (m *MockQuizService) GetByCreator(id uuid.UUID) ([]*domain.Quiz, error) {
	m.getByCreatorCalledWith = id
	return m.getByCreatorReturns, m.getByCreatorReturnsError
}

func (m *MockQuizService) GetByID(uuid.UUID) (*domain.Quiz, error) {
	return m.getByIdReturns, m.getByIdReturnsError
}

func (m *MockQuizService) CreateOrUpdate(quiz *domain.Quiz) error {
	m.createOrUpdateCalledWith = quiz
	return m.createOrUpdateReturns
}
func (m *MockQuizService) Delete(id uuid.UUID) error {
	m.deleteCalledWith = id
	return m.deleteReturns
}

type MockJwtService struct {
	services.JwtService

	validateTokenCalledWith   string
	validateTokenReturns      *jwt.Token
	validateTokenReturnsError error
}

func (m *MockJwtService) ValidateToken(token string) (*jwt.Token, error) {
	m.validateTokenCalledWith = token
	return m.validateTokenReturns, m.validateTokenReturnsError
}
