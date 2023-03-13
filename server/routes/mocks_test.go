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

type MockPlayerService struct {
	services.PlayerService

	getByGameCalledWith   uuid.UUID
	getByGameReturns      []*domain.Player
	getByGameReturnsError error

	createCalledWith *domain.Player
	createReturns    error

	getByIdReturns      *domain.Player
	getByIdReturnsError error

	deleteReturns error
}

func (m *MockPlayerService) GetByGame(id uuid.UUID) ([]*domain.Player, error) {
	m.getByGameCalledWith = id
	return m.getByGameReturns, m.getByGameReturnsError
}
func (m *MockPlayerService) GetByID(uuid.UUID) (*domain.Player, error) {
	return m.getByIdReturns, m.getByIdReturnsError
}

func (m *MockPlayerService) Delete(*domain.Player) error {
	return m.deleteReturns
}

func (m *MockPlayerService) Create(player *domain.Player) error {
	m.createCalledWith = player
	return m.createReturns
}

type MockGameService struct {
	services.GameService

	getByQuizCalledWith   uuid.UUID
	getByQuizReturns      []*domain.Game
	getByQuizReturnsError error

	createCalledWith *domain.Game
	createReturns    error

	getByIdReturns      *domain.Game
	getByIdReturnsError error

	startCalledWith *domain.Game
	startReturns    error

	finishCalledWith *domain.Game
	finishReturns    error

	deleteCalledWith *domain.Game
	deleteReturns    error
}

func (m *MockGameService) GetByQuiz(id uuid.UUID) ([]*domain.Game, error) {
	m.getByQuizCalledWith = id
	return m.getByQuizReturns, m.getByQuizReturnsError
}

func (m *MockGameService) GetByID(_ uuid.UUID) (*domain.Game, error) {
	return m.getByIdReturns, m.getByIdReturnsError
}

func (m *MockGameService) Create(game *domain.Game) error {
	m.createCalledWith = game
	return m.createReturns
}

func (m *MockGameService) Delete(game *domain.Game) error {
	m.deleteCalledWith = game
	return m.deleteReturns
}

func (m *MockGameService) Start(game *domain.Game) error {
	m.startCalledWith = game
	return m.startReturns
}

func (m *MockGameService) Finish(game *domain.Game) error {
	m.finishCalledWith = game
	return m.finishReturns
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
