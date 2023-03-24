package routes

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/survivorbat/qq.maarten.dev/server/coordinator"
	"github.com/survivorbat/qq.maarten.dev/server/domain"
	"github.com/survivorbat/qq.maarten.dev/server/services"
	"sync"
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

	nextCalledWith *domain.Game
	nextReturns    error

	answerReturns error

	getByCodeCalledWith   string
	getByCodeReturns      *domain.Game
	getByCodeReturnsError error
}

func (m *MockGameService) GetByQuiz(id uuid.UUID) ([]*domain.Game, error) {
	m.getByQuizCalledWith = id
	return m.getByQuizReturns, m.getByQuizReturnsError
}

func (m *MockGameService) GetByCode(code string) (*domain.Game, error) {
	m.getByCodeCalledWith = code
	return m.getByCodeReturns, m.getByCodeReturnsError
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

func (m *MockGameService) Next(game *domain.Game) error {
	m.nextCalledWith = game
	return m.nextReturns
}

func (m *MockGameService) AnswerQuestion(*domain.Game, uuid.UUID, uuid.UUID, uuid.UUID) error {
	return m.answerReturns
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

type MockCoordinator struct {
	coordinator.GameCoordinator

	subscribeCreatorCallbackCalledWithGame    uuid.UUID
	subscribeCreatorCallbackCalledWithCreator *domain.Creator
	subscribeCreatorCallbackReturns           *coordinator.BroadcastMessage

	handleCreatorMessageWaitGroup         sync.WaitGroup
	handleCreatorMessagePanicsWith        any
	handleCreatorMessageCalledWithGame    uuid.UUID
	handleCreatorMessageCalledWithMessage *coordinator.CreatorMessage

	unsubscribeCreatorWaitGroup      sync.WaitGroup
	unsubscribeCreatorCalledWithGame uuid.UUID

	subscribePlayerCallbackCalledWithGame   uuid.UUID
	subscribePlayerCallbackCalledWithPlayer *domain.Player
	subscribePlayerCallbackReturns          *coordinator.BroadcastMessage

	unsubscribePlayerWaitGroup        sync.WaitGroup
	unsubscribePlayerCalledWithGame   uuid.UUID
	unsubscribePlayerCalledWithPlayer *domain.Player

	handlePlayerMessageWaitGroup         sync.WaitGroup
	handlePlayerMessageCalledWithGame    uuid.UUID
	handlePlayerMessageCalledWithPlayer  uuid.UUID
	handlePlayerMessageCalledWithMessage *coordinator.PlayerMessage
	handlePlayerMessagePanicsWith        any
}

func (m *MockCoordinator) SubscribeCreator(gameID uuid.UUID, creator *domain.Creator, callback coordinator.BroadcastCallback) {
	m.subscribeCreatorCallbackCalledWithGame = gameID
	m.subscribeCreatorCallbackCalledWithCreator = creator

	if m.subscribeCreatorCallbackReturns != nil {
		callback(m.subscribeCreatorCallbackReturns)
	}
}

func (m *MockCoordinator) SubscribePlayer(gameID uuid.UUID, player *domain.Player, callback coordinator.BroadcastCallback) {
	m.subscribePlayerCallbackCalledWithGame = gameID
	m.subscribePlayerCallbackCalledWithPlayer = player

	if m.subscribePlayerCallbackReturns != nil {
		callback(m.subscribePlayerCallbackReturns)
	}
}

func (m *MockCoordinator) UnsubscribeCreator(gameId uuid.UUID) {
	defer m.unsubscribeCreatorWaitGroup.Done()
	m.unsubscribeCreatorCalledWithGame = gameId
}

func (m *MockCoordinator) UnsubscribePlayer(gameId uuid.UUID, player *domain.Player) {
	defer m.unsubscribePlayerWaitGroup.Done()
	m.unsubscribePlayerCalledWithGame = gameId
	m.unsubscribePlayerCalledWithPlayer = player
}

func (m *MockCoordinator) HandleCreatorMessage(gameID uuid.UUID, action *coordinator.CreatorMessage) {
	defer m.handleCreatorMessageWaitGroup.Done()

	m.handleCreatorMessageCalledWithGame = gameID
	m.handleCreatorMessageCalledWithMessage = action

	if m.handleCreatorMessagePanicsWith != nil {
		panic(m.handleCreatorMessagePanicsWith)
	}
}

func (m *MockCoordinator) HandlePlayerMessage(gameID uuid.UUID, playerID uuid.UUID, action *coordinator.PlayerMessage) {
	defer m.handlePlayerMessageWaitGroup.Done()

	m.handlePlayerMessageCalledWithGame = gameID
	m.handlePlayerMessageCalledWithPlayer = playerID
	m.handlePlayerMessageCalledWithMessage = action

	if m.handlePlayerMessagePanicsWith != nil {
		panic(m.handlePlayerMessagePanicsWith)
	}
}
