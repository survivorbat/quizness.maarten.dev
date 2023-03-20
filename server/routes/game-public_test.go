package routes

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/survivorbat/qq.maarten.dev/server/domain"
	"github.com/survivorbat/qq.maarten.dev/server/routes/outputs"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPublicGameHandler_GetByCode_ReturnsErrorOnGameNotFound(t *testing.T) {
	t.Parallel()
	// Arrange
	gameService := &MockGameService{getByCodeReturnsError: assert.AnError}
	handler := &PublicGameHandler{GameService: gameService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Request, _ = http.NewRequest(http.MethodGet, "?code=abc", nil)

	// Act
	handler.GetByCode(context)

	// Assert
	assert.Equal(t, http.StatusNotFound, writer.Code)
}

func TestPublicGameHandler_GetByCode_ReturnsErrorOnNoCode(t *testing.T) {
	t.Parallel()
	// Arrange
	gameService := &MockGameService{getByCodeReturnsError: assert.AnError}
	handler := &PublicGameHandler{GameService: gameService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Request, _ = http.NewRequest(http.MethodGet, "", nil)

	// Act
	handler.GetByCode(context)

	// Assert
	assert.Equal(t, http.StatusForbidden, writer.Code)
}

func TestPublicGameHandler_GetByCode_ReturnsGameID(t *testing.T) {
	t.Parallel()
	// Arrange
	game := &domain.Game{BaseObject: domain.BaseObject{ID: uuid.MustParse("788f12a9-51e8-4c87-9b0c-06bcc9f0691b")}}
	gameService := &MockGameService{getByCodeReturns: game}
	handler := &PublicGameHandler{GameService: gameService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Request, _ = http.NewRequest(http.MethodGet, "?code=abc", nil)

	// Act
	handler.GetByCode(context)

	// Assert
	assert.Equal(t, http.StatusOK, writer.Code)

	body, err := io.ReadAll(writer.Body)
	if err != nil {
		t.Fatal(err)
	}

	var result *outputs.OutputGame
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, game.ID, result.ID)
}

func TestPublicGameHandler_GetQuiz_ReturnsErrorOnInvalidUUID(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := &PublicGameHandler{}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	context.Params = []gin.Param{{Key: "id", Value: "no"}}

	// Act
	handler.GetQuiz(context)

	// Assert
	assert.Equal(t, http.StatusBadRequest, writer.Code)
}

func TestPublicGameHandler_GetQuiz_ReturnsErrorOnGameNotFound(t *testing.T) {
	t.Parallel()
	// Arrange
	gameService := &MockGameService{getByIdReturnsError: assert.AnError}
	handler := &PublicGameHandler{GameService: gameService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	context.Params = []gin.Param{{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"}}

	// Act
	handler.GetQuiz(context)

	// Assert
	assert.Equal(t, http.StatusNotFound, writer.Code)
}

func TestPublicGameHandler_GetByQuiz_ReturnsQuiz(t *testing.T) {
	t.Parallel()
	// Arrange
	game := &domain.Game{
		BaseObject: domain.BaseObject{ID: uuid.MustParse("788f12a9-51e8-4c87-9b0c-06bcc9f0691b")},
		Quiz: &domain.Quiz{
			Name: "abc",
		},
	}
	gameService := &MockGameService{getByIdReturns: game}
	handler := &PublicGameHandler{GameService: gameService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	context.Params = []gin.Param{{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"}}
	// Act
	handler.GetQuiz(context)

	// Assert
	assert.Equal(t, http.StatusOK, writer.Code)

	body, err := io.ReadAll(writer.Body)
	if err != nil {
		t.Fatal(err)
	}

	var result *outputs.OutputQuiz
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, game.Quiz.Name, result.Name)
}
