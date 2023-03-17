package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/survivorbat/qq.maarten.dev/server/domain"
	"net/http"
	"net/http/httptest"
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
