package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/survivorbat/qq.maarten.dev/server/coordinator"
	"github.com/survivorbat/qq.maarten.dev/server/domain"
	"github.com/survivorbat/qq.maarten.dev/server/inputs"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGameConnectionHandler_Get_ReturnsErrorOnInvalidGameUUID(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := &GameConnectionHandler{}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	context.Params = []gin.Param{{Key: "id", Value: "no"}}

	// Act
	handler.Get(context)

	// Assert
	assert.Equal(t, http.StatusBadRequest, writer.Code)
}

func TestGameConnectionHandler_Get_ReturnsErrorOnInvalidPlayerUUID(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := &GameConnectionHandler{}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	context.Params = []gin.Param{
		{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"},
		{Key: "player", Value: "no"},
	}

	// Act
	handler.Get(context)

	// Assert
	assert.Equal(t, http.StatusBadRequest, writer.Code)
}

func TestGameConnectionHandler_Get_ReturnsErrorOnGameNotFound(t *testing.T) {
	t.Parallel()
	// Arrange
	gameService := &MockGameService{getByIdReturnsError: assert.AnError}
	handler := &GameConnectionHandler{GameService: gameService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	context.Params = []gin.Param{
		{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"},
		{Key: "player", Value: "3ad4afb5-91af-4243-b06f-40089db9a63a"},
	}

	// Act
	handler.Get(context)

	// Assert
	assert.Equal(t, http.StatusNotFound, writer.Code)
}

func TestGameConnectionHandler_Get_ReturnsErrorOnPlayerNotInGame(t *testing.T) {
	t.Parallel()
	// Arrange
	game := &domain.Game{}
	gameService := &MockGameService{getByIdReturns: game}
	handler := &GameConnectionHandler{GameService: gameService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	context.Params = []gin.Param{
		{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"},
		{Key: "player", Value: "3ad4afb5-91af-4243-b06f-40089db9a63a"},
	}

	// Act
	handler.Get(context)

	// Assert
	assert.Equal(t, http.StatusForbidden, writer.Code)
}

func TestGameConnectionHandler_Get_ReturnsErrorOnPlayerNotFound(t *testing.T) {
	t.Parallel()
	// Arrange
	playerID := uuid.MustParse("3ad4afb5-91af-4243-b06f-40089db9a63a")
	game := &domain.Game{Players: []*domain.Player{{BaseObject: domain.BaseObject{ID: playerID}}}}

	playerService := &MockPlayerService{getByIdReturnsError: assert.AnError}
	gameService := &MockGameService{getByIdReturns: game}
	handler := &GameConnectionHandler{GameService: gameService, PlayerService: playerService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	context.Params = []gin.Param{
		{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"},
		{Key: "player", Value: playerID.String()},
	}

	// Act
	handler.Get(context)

	// Assert
	assert.Equal(t, http.StatusNotFound, writer.Code)
}
func TestGameConnectionHandler_Get_ReturnsErrorOnSocketHeadersNotFound(t *testing.T) {
	t.Parallel()
	// Arrange
	playerID := uuid.MustParse("3ad4afb5-91af-4243-b06f-40089db9a63a")
	game := &domain.Game{Players: []*domain.Player{{BaseObject: domain.BaseObject{ID: playerID}}}}

	playerService := &MockPlayerService{}
	gameService := &MockGameService{getByIdReturns: game}
	handler := &GameConnectionHandler{GameService: gameService, PlayerService: playerService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	context.Params = []gin.Param{
		{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"},
		{Key: "player", Value: playerID.String()},
	}

	// Act
	handler.Get(context)

	// Assert
	assert.Equal(t, http.StatusBadRequest, writer.Code)
}

func TestGameConnectionHandler_GetCreator_ReturnsErrorOnInvalidGameUUID(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := &GameConnectionHandler{}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Params = []gin.Param{{Key: "id", Value: "no"}}

	// Act
	handler.GetCreator(context)

	// Assert
	assert.Equal(t, http.StatusBadRequest, writer.Code)
}

func TestGameConnectionHandler_GetCreator_ReturnsErrorOnGameNotFound(t *testing.T) {
	t.Parallel()
	// Arrange
	gameService := &MockGameService{getByIdReturnsError: assert.AnError}
	handler := &GameConnectionHandler{GameService: gameService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	context.Params = []gin.Param{
		{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"},
	}

	// Act
	handler.GetCreator(context)

	// Assert
	assert.Equal(t, http.StatusNotFound, writer.Code)
}

func TestGameConnectionHandler_GetCreator_ReturnsErrorOnDifferentCreator(t *testing.T) {
	t.Parallel()
	// Arrange
	game := &domain.Game{Quiz: &domain.Quiz{CreatorID: uuid.MustParse("31251895-88a3-4111-86fd-f5291bfb9c69")}}
	gameService := &MockGameService{getByIdReturns: game}
	handler := &GameConnectionHandler{GameService: gameService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	context.Params = []gin.Param{
		{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"},
	}

	// Act
	handler.GetCreator(context)

	// Assert
	assert.Equal(t, http.StatusForbidden, writer.Code)
}

func TestGameConnectionHandler_GetCreator_ReturnsErrorOnCreatorNotFound(t *testing.T) {
	t.Parallel()
	// Arrange
	game := &domain.Game{
		Quiz: &domain.Quiz{
			CreatorID: uuid.MustParse("2f80947c-e724-4b38-8c8d-3823864fef58"),
		},
	}

	creatorService := &MockCreatorService{getByIDReturnsError: assert.AnError}
	gameService := &MockGameService{getByIdReturns: game}
	handler := &GameConnectionHandler{GameService: gameService, CreatorService: creatorService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	context.Params = []gin.Param{
		{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"},
	}

	// Act
	handler.GetCreator(context)

	// Assert
	assert.Equal(t, http.StatusNotFound, writer.Code)
}

func TestGameConnectionHandler_GetCreator_ReturnsErrorOnSocketWrongHeaders(t *testing.T) {
	t.Parallel()
	// Arrange
	game := &domain.Game{
		Quiz: &domain.Quiz{
			CreatorID: uuid.MustParse("2f80947c-e724-4b38-8c8d-3823864fef58"),
		},
	}

	creatorService := &MockCreatorService{getByIDReturns: &domain.Creator{}}
	gameService := &MockGameService{getByIdReturns: game}
	handler := &GameConnectionHandler{GameService: gameService, CreatorService: creatorService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	context.Params = []gin.Param{
		{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"},
	}

	// Act
	handler.GetCreator(context)

	// Assert
	assert.Equal(t, http.StatusBadRequest, writer.Code)
}

// Websocket tests

func TestGameConnectionHandler_Get_WritesBroadcastsToSocket(t *testing.T) {
	t.Parallel()
	// Arrange
	playerID := uuid.MustParse("3ad4afb5-91af-4243-b06f-40089db9a63a")
	game := &domain.Game{
		BaseObject: domain.BaseObject{ID: uuid.MustParse("f7422157-bc0c-4998-834a-0aeb7a800dc7")},
		Players:    []*domain.Player{{BaseObject: domain.BaseObject{ID: playerID}}},
	}

	coord := &MockCoordinator{
		subscribePlayerCallbackReturns: &coordinator.BroadcastMessage{
			Type: coordinator.FinishGameType,
		},
	}
	playerService := &MockPlayerService{getByIdReturns: game.Players[0]}
	gameService := &MockGameService{getByIdReturns: game}
	handler := &GameConnectionHandler{GameService: gameService, PlayerService: playerService, Coordinator: coord}

	engine := gin.Default()
	engine.GET("/games/:id/players/:player/connection", handler.Get)
	ts := httptest.NewServer(engine)
	defer ts.Close()

	socketUrl := fmt.Sprintf("ws%s/games/f7422157-bc0c-4998-834a-0aeb7a800dc7/players/3ad4afb5-91af-4243-b06f-40089db9a63a/connection", strings.TrimPrefix(ts.URL, "http"))

	// Act
	ws, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)

	// Assert
	if !assert.NoError(t, err) {
		t.Fatal(err)
	}
	defer ws.Close()

	var message coordinator.BroadcastMessage
	if err := ws.ReadJSON(&message); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, coordinator.FinishGameType, message.Type)
	assert.Equal(t, playerService.getByIdReturns, coord.subscribePlayerCallbackCalledWithPlayer)
	assert.Equal(t, game.ID, coord.subscribePlayerCallbackCalledWithGame)
}

func TestGameConnectionHandler_Get_SendsMessagesFromSocket(t *testing.T) {
	t.Parallel()
	// Arrange
	playerID := uuid.MustParse("3ad4afb5-91af-4243-b06f-40089db9a63a")
	game := &domain.Game{
		BaseObject: domain.BaseObject{ID: uuid.MustParse("f7422157-bc0c-4998-834a-0aeb7a800dc7")},
		Players:    []*domain.Player{{BaseObject: domain.BaseObject{ID: playerID}}},
	}

	coord := &MockCoordinator{}
	coord.handlePlayerMessageWaitGroup.Add(1)
	playerService := &MockPlayerService{getByIdReturns: game.Players[0]}
	gameService := &MockGameService{getByIdReturns: game}
	handler := &GameConnectionHandler{GameService: gameService, PlayerService: playerService, Coordinator: coord}

	engine := gin.Default()
	engine.GET("/games/:id/players/:player/connection", handler.Get)
	ts := httptest.NewServer(engine)
	defer ts.Close()

	socketUrl := fmt.Sprintf("ws%s/games/f7422157-bc0c-4998-834a-0aeb7a800dc7/players/3ad4afb5-91af-4243-b06f-40089db9a63a/connection", strings.TrimPrefix(ts.URL, "http"))

	// Act
	ws, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)

	// Assert
	if !assert.NoError(t, err) {
		t.Fatal(err)
	}
	defer ws.Close()

	answer := &inputs.Answer{OptionID: uuid.MustParse("1ca32a64-3f5c-4bc8-9c13-0b3d2f1600bc")}
	answerJson, _ := json.Marshal(answer)
	action := &coordinator.PlayerMessage{Action: coordinator.AnswerAction, Content: answerJson}
	if err := ws.WriteJSON(&action); err != nil {
		t.Fatal(err)
	}

	coord.handlePlayerMessageWaitGroup.Wait()

	assert.Equal(t, game.ID, coord.handlePlayerMessageCalledWithGame)
	assert.Equal(t, answer, coord.handlePlayerMessageCalledWithMessage.Answer)
}
func TestGameConnectionHandler_Get_CallsUnsubscribeOnPanic(t *testing.T) {
	t.Parallel()
	// Arrange
	playerID := uuid.MustParse("3ad4afb5-91af-4243-b06f-40089db9a63a")
	game := &domain.Game{
		BaseObject: domain.BaseObject{ID: uuid.MustParse("f7422157-bc0c-4998-834a-0aeb7a800dc7")},
		Players:    []*domain.Player{{BaseObject: domain.BaseObject{ID: playerID}}},
	}

	coord := &MockCoordinator{handlePlayerMessagePanicsWith: "abc"}
	coord.handlePlayerMessageWaitGroup.Add(1)
	playerService := &MockPlayerService{getByIdReturns: game.Players[0]}
	gameService := &MockGameService{getByIdReturns: game}
	handler := &GameConnectionHandler{GameService: gameService, PlayerService: playerService, Coordinator: coord}

	engine := gin.Default()
	engine.GET("/games/:id/players/:player/connection", handler.Get)
	ts := httptest.NewServer(engine)
	defer ts.Close()

	socketUrl := fmt.Sprintf("ws%s/games/f7422157-bc0c-4998-834a-0aeb7a800dc7/players/3ad4afb5-91af-4243-b06f-40089db9a63a/connection", strings.TrimPrefix(ts.URL, "http"))

	// Act
	ws, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)

	// Assert
	if !assert.NoError(t, err) {
		t.Fatal(err)
	}
	defer ws.Close()

	answer := &inputs.Answer{OptionID: uuid.MustParse("1ca32a64-3f5c-4bc8-9c13-0b3d2f1600bc")}
	answerJson, _ := json.Marshal(answer)
	action := &coordinator.PlayerMessage{Action: coordinator.AnswerAction, Content: answerJson}
	if err := ws.WriteJSON(&action); err != nil {
		t.Fatal(err)
	}

	coord.handlePlayerMessageWaitGroup.Wait()

	assert.Equal(t, game.ID, coord.unsubscribePlayerCalledWithGame)
}

func TestGameConnectionHandler_GetCreator_WritesBroadcastsToSocket(t *testing.T) {
	t.Parallel()
	// Arrange
	game := &domain.Game{
		BaseObject: domain.BaseObject{ID: uuid.MustParse("f7422157-bc0c-4998-834a-0aeb7a800dc7")},
		Quiz: &domain.Quiz{
			CreatorID: uuid.MustParse("2f80947c-e724-4b38-8c8d-3823864fef58"),
		},
	}

	coord := &MockCoordinator{
		subscribeCreatorCallbackReturns: &coordinator.BroadcastMessage{
			Type: coordinator.FinishGameType,
		},
	}
	creatorService := &MockCreatorService{getByIDReturns: &domain.Creator{Nickname: "abc"}}
	gameService := &MockGameService{getByIdReturns: game}
	handler := &GameConnectionHandler{GameService: gameService, CreatorService: creatorService, Coordinator: coord}

	engine := gin.Default()

	// Middleware the user in
	engine.Use(func(context *gin.Context) {
		context.Set("user", game.Quiz.CreatorID.String())
	})

	engine.GET("/games/:id/connection", handler.GetCreator)
	ts := httptest.NewServer(engine)
	defer ts.Close()

	socketUrl := fmt.Sprintf("ws%s/games/f7422157-bc0c-4998-834a-0aeb7a800dc7/connection", strings.TrimPrefix(ts.URL, "http"))

	// Act
	ws, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)

	// Assert
	if !assert.NoError(t, err) {
		t.Fatal(err)
	}
	defer ws.Close()

	var message coordinator.BroadcastMessage
	if err := ws.ReadJSON(&message); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, coordinator.FinishGameType, message.Type)
	assert.Equal(t, creatorService.getByIDReturns, coord.subscribeCreatorCallbackCalledWithCreator)
	assert.Equal(t, game.ID, coord.subscribeCreatorCallbackCalledWithGame)
}

func TestGameConnectionHandler_GetCreator_SendsMessagesFromSocket(t *testing.T) {
	t.Parallel()
	// Arrange
	game := &domain.Game{
		BaseObject: domain.BaseObject{ID: uuid.MustParse("f7422157-bc0c-4998-834a-0aeb7a800dc7")},
		Quiz: &domain.Quiz{
			CreatorID: uuid.MustParse("2f80947c-e724-4b38-8c8d-3823864fef58"),
		},
	}

	coord := &MockCoordinator{}
	coord.handleCreatorMessageWaitGroup.Add(1)

	creatorService := &MockCreatorService{getByIDReturns: &domain.Creator{}}
	gameService := &MockGameService{getByIdReturns: game}
	handler := &GameConnectionHandler{GameService: gameService, CreatorService: creatorService, Coordinator: coord}

	engine := gin.Default()

	// Middleware the user in
	engine.Use(func(context *gin.Context) {
		context.Set("user", game.Quiz.CreatorID.String())
	})

	engine.GET("/games/:id/connection", handler.GetCreator)
	ts := httptest.NewServer(engine)
	defer ts.Close()

	socketUrl := fmt.Sprintf("ws%s/games/f7422157-bc0c-4998-834a-0aeb7a800dc7/connection", strings.TrimPrefix(ts.URL, "http"))

	// Act
	ws, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)

	// Assert
	if !assert.NoError(t, err) {
		t.Fatal(err)
	}
	defer ws.Close()

	action := &coordinator.CreatorMessage{Action: coordinator.FinishGameAction}
	if err := ws.WriteJSON(&action); err != nil {
		t.Fatal(err)
	}

	coord.handleCreatorMessageWaitGroup.Wait()

	assert.Equal(t, game.ID, coord.handleCreatorMessageCalledWithGame)
	assert.Equal(t, action, coord.handleCreatorMessageCalledWithMessage)
}

func TestGameConnectionHandler_GetCreator_PanicCallsUnsubscribe(t *testing.T) {
	t.Parallel()
	// Arrange
	game := &domain.Game{
		BaseObject: domain.BaseObject{ID: uuid.MustParse("f7422157-bc0c-4998-834a-0aeb7a800dc7")},
		Quiz: &domain.Quiz{
			CreatorID: uuid.MustParse("2f80947c-e724-4b38-8c8d-3823864fef58"),
		},
	}

	coord := &MockCoordinator{handleCreatorMessagePanicsWith: "test"}
	coord.handleCreatorMessageWaitGroup.Add(1)

	creatorService := &MockCreatorService{getByIDReturns: &domain.Creator{}}
	gameService := &MockGameService{getByIdReturns: game}
	handler := &GameConnectionHandler{GameService: gameService, CreatorService: creatorService, Coordinator: coord}

	engine := gin.Default()

	// Middleware the user in
	engine.Use(func(context *gin.Context) {
		context.Set("user", game.Quiz.CreatorID.String())
	})

	engine.GET("/games/:id/connection", handler.GetCreator)
	ts := httptest.NewServer(engine)
	defer ts.Close()

	socketUrl := fmt.Sprintf("ws%s/games/f7422157-bc0c-4998-834a-0aeb7a800dc7/connection", strings.TrimPrefix(ts.URL, "http"))

	// Act
	ws, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)

	// Assert
	if !assert.NoError(t, err) {
		t.Fatal(err)
	}
	defer ws.Close()

	action := &coordinator.CreatorMessage{Action: coordinator.FinishGameAction}
	if err := ws.WriteJSON(&action); err != nil {
		t.Fatal(err)
	}

	coord.handleCreatorMessageWaitGroup.Wait()

	assert.Equal(t, game.ID, coord.unsubscribeCreatorCalledWithGame)
}
