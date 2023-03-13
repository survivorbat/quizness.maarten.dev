package routes

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/survivorbat/qq.maarten.dev/server/domain"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPlayerHandler_Get_ReturnsErrorOnInvalidUUID(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := &PlayerHandler{}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	context.Params = []gin.Param{{Key: "id", Value: "no"}}

	// Act
	handler.Get(context)

	// Assert
	assert.Equal(t, http.StatusBadRequest, writer.Code)
}

func TestPlayerHandler_Get_ReturnsErrorOnGameNotFound(t *testing.T) {
	t.Parallel()
	// Arrange
	gameService := &MockGameService{getByIdReturnsError: assert.AnError}
	handler := &PlayerHandler{GameService: gameService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	context.Params = []gin.Param{{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"}}

	// Act
	handler.Get(context)

	// Assert
	assert.Equal(t, http.StatusNotFound, writer.Code)
}

func TestPlayerHandler_Get_ReturnsErrorOnNotMyGame(t *testing.T) {
	t.Parallel()
	// Arrange
	gameService := &MockGameService{getByIdReturns: &domain.Game{
		Quiz: &domain.Quiz{
			CreatorID: uuid.MustParse("8fdc3e5a-b0a8-4103-af3b-c2f20d91889b"),
		},
	}}
	handler := &PlayerHandler{GameService: gameService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	context.Params = []gin.Param{{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"}}

	// Act
	handler.Get(context)

	// Assert
	assert.Equal(t, http.StatusForbidden, writer.Code)
}

func TestPlayerHandler_Get_ReturnsErrorOnFailedGet(t *testing.T) {
	t.Parallel()
	// Arrange
	gameService := &MockGameService{getByIdReturns: &domain.Game{
		Quiz: &domain.Quiz{
			CreatorID: uuid.MustParse("2f80947c-e724-4b38-8c8d-3823864fef58"),
		},
	}}
	playerService := &MockPlayerService{getByGameReturnsError: assert.AnError}
	handler := &PlayerHandler{GameService: gameService, PlayerService: playerService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	context.Params = []gin.Param{{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"}}

	// Act
	handler.Get(context)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, writer.Code)
}

func TestPlayerHandler_Get_ReturnsData(t *testing.T) {
	t.Parallel()
	// Arrange
	gameService := &MockGameService{getByIdReturns: &domain.Game{
		Quiz: &domain.Quiz{
			CreatorID: uuid.MustParse("2f80947c-e724-4b38-8c8d-3823864fef58"),
		},
	}}
	playerService := &MockPlayerService{getByGameReturns: []*domain.Player{{Nickname: "A"}, {Nickname: "B"}}}
	handler := &PlayerHandler{GameService: gameService, PlayerService: playerService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	context.Params = []gin.Param{{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"}}

	// Act
	handler.Get(context)

	// Assert
	assert.Equal(t, http.StatusOK, writer.Code)

	var result []*domain.Player
	if err := json.Unmarshal(writer.Body.Bytes(), &result); err != nil {
		t.Fatal(err.Error())
	}

	assert.ElementsMatch(t, playerService.getByGameReturns, result)
}

func TestPlayerHandler_Post_ReturnsErrorOnInvalidUUID(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := &PlayerHandler{}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Request, _ = http.NewRequest(http.MethodPost, "", nil)
	context.Params = []gin.Param{{Key: "id", Value: "no"}}

	// Act
	handler.Post(context)

	// Assert
	assert.Equal(t, http.StatusBadRequest, writer.Code)
}

func TestPlayerHandler_Post_ReturnsErrorOnGameNotFound(t *testing.T) {
	t.Parallel()
	// Arrange
	gameService := &MockGameService{getByIdReturnsError: assert.AnError}
	handler := &PlayerHandler{GameService: gameService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Request, _ = http.NewRequest(http.MethodPost, "", nil)
	context.Params = []gin.Param{{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"}}

	// Act
	handler.Post(context)

	// Assert
	assert.Equal(t, http.StatusNotFound, writer.Code)
}

func TestPlayerHandler_Post_ReturnsErrorOnNotInProgress(t *testing.T) {
	t.Parallel()
	// Arrange
	gameService := &MockGameService{getByIdReturns: &domain.Game{}}
	handler := &PlayerHandler{GameService: gameService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Request, _ = http.NewRequest(http.MethodPost, "", nil)
	context.Params = []gin.Param{{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"}}

	// Act
	handler.Post(context)

	// Assert
	assert.Equal(t, http.StatusNotFound, writer.Code)
}

func TestPlayerHandler_Post_ReturnsErrorOnFull(t *testing.T) {
	t.Parallel()
	// Arrange
	gameService := &MockGameService{getByIdReturns: &domain.Game{StartTime: time.Now()}}
	handler := &PlayerHandler{GameService: gameService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Request, _ = http.NewRequest(http.MethodPost, "", nil)
	context.Params = []gin.Param{{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"}}

	// Act
	handler.Post(context)

	// Assert
	assert.Equal(t, http.StatusConflict, writer.Code)
}

func TestPlayerHandler_Post_ReturnsCreateError(t *testing.T) {
	t.Parallel()
	// Arrange
	gameService := &MockGameService{getByIdReturns: &domain.Game{StartTime: time.Now(), PlayerLimit: 5}}
	playerService := &MockPlayerService{createReturns: assert.AnError}
	handler := &PlayerHandler{GameService: gameService, PlayerService: playerService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Request, _ = http.NewRequest(http.MethodPost, "", nil)
	context.Params = []gin.Param{{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"}}

	// Act
	handler.Post(context)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, writer.Code)
}

func TestPlayerHandler_Post_ReturnsNewPlayer(t *testing.T) {
	t.Parallel()
	// Arrange
	gameService := &MockGameService{getByIdReturns: &domain.Game{StartTime: time.Now(), PlayerLimit: 5}}
	playerService := &MockPlayerService{}
	handler := &PlayerHandler{GameService: gameService, PlayerService: playerService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Request, _ = http.NewRequest(http.MethodPost, "", nil)
	context.Params = []gin.Param{{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"}}

	// Act
	handler.Post(context)

	// Assert
	assert.Equal(t, http.StatusOK, writer.Code)

	var result *domain.Player
	if err := json.Unmarshal(writer.Body.Bytes(), &result); err != nil {
		t.Fatal(err.Error())
	}

	assert.Equal(t, gameService.getByIdReturns.ID, result.GameID)
}

func TestPlayerHandler_Delete_ReturnsErrorOnInvalidUUID(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := &PlayerHandler{}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Request, _ = http.NewRequest(http.MethodDelete, "", nil)
	context.Params = []gin.Param{{Key: "id", Value: "no"}}

	// Act
	handler.Delete(context)

	// Assert
	assert.Equal(t, http.StatusBadRequest, writer.Code)
}

func TestPlayerHandler_Delete_ReturnsErrorOnGameNotFound(t *testing.T) {
	t.Parallel()
	// Arrange
	playerService := &MockPlayerService{getByIdReturnsError: assert.AnError}
	handler := &PlayerHandler{PlayerService: playerService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Request, _ = http.NewRequest(http.MethodDelete, "", nil)
	context.Params = []gin.Param{{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"}}

	// Act
	handler.Delete(context)

	// Assert
	assert.Equal(t, http.StatusNotFound, writer.Code)
}

func TestPlayerHandler_Delete_ReturnsDeleteError(t *testing.T) {
	t.Parallel()
	// Arrange
	playerService := &MockPlayerService{deleteReturns: assert.AnError}
	handler := &PlayerHandler{PlayerService: playerService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Request, _ = http.NewRequest(http.MethodDelete, "", nil)
	context.Params = []gin.Param{{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"}}

	// Act
	handler.Delete(context)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, writer.Code)
}

func TestPlayerHandler_Delete_ReturnsSuccess(t *testing.T) {
	t.Parallel()
	// Arrange
	playerService := &MockPlayerService{}
	handler := &PlayerHandler{PlayerService: playerService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Request, _ = http.NewRequest(http.MethodDelete, "", nil)
	context.Params = []gin.Param{{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"}}

	// Act
	handler.Delete(context)

	// Assert
	assert.Equal(t, http.StatusOK, writer.Code)
}
