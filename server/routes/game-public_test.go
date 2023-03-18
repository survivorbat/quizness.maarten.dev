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
