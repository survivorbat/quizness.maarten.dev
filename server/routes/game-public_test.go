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

func TestPublicGameHandler_Get_ReturnsErrorOnInvalidGameUUID(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := &PublicGameHandler{}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	context.Params = []gin.Param{{Key: "id", Value: "no"}}

	// Act
	handler.Get(context)

	// Assert
	assert.Equal(t, http.StatusBadRequest, writer.Code)
}

func TestPublicGameHandler_Get_ReturnsErrorOnGameNotFound(t *testing.T) {
	t.Parallel()
	// Arrange
	gameService := &MockGameService{getByIdReturnsError: assert.AnError}
	handler := &PublicGameHandler{GameService: gameService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	context.Params = []gin.Param{
		{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"},
	}

	// Act
	handler.Get(context)

	// Assert
	assert.Equal(t, http.StatusNotFound, writer.Code)
}

func TestPublicGameHandler_Get_ReturnsErrorOnQuestionNotFound(t *testing.T) {
	t.Parallel()
	// Arrange
	game := &domain.Game{Quiz: &domain.Quiz{}}
	gameService := &MockGameService{getByIdReturns: game}
	handler := &PublicGameHandler{GameService: gameService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	context.Params = []gin.Param{
		{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"},
	}

	// Act
	handler.Get(context)

	// Assert
	assert.Equal(t, http.StatusNotFound, writer.Code)
}

func TestPublicGameHandler_Get_ReturnsQuestion(t *testing.T) {
	t.Parallel()
	// Arrange
	question := &domain.MultipleChoiceQuestion{
		BaseQuestion: domain.BaseQuestion{BaseObject: domain.BaseObject{ID: uuid.MustParse("78fdb37a-075a-4ed3-ae5f-ffa7d7bca781")}},
	}
	game := &domain.Game{
		CurrentQuestion: question.ID,
		Quiz: &domain.Quiz{
			MultipleChoiceQuestions: []*domain.MultipleChoiceQuestion{question},
		},
	}
	gameService := &MockGameService{getByIdReturns: game}
	handler := &PublicGameHandler{GameService: gameService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	context.Params = []gin.Param{
		{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"},
	}

	// Act
	handler.Get(context)

	// Assert
	assert.Equal(t, http.StatusOK, writer.Code)

	body, err := io.ReadAll(writer.Body)
	if err != nil {
		t.Fatal(err)
	}

	var result *outputs.OutputMultipleChoiceQuestion
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, question.ID, result.ID)
	assert.Equal(t, question.Title, result.Title)
	assert.Equal(t, question.Description, result.Description)
	assert.Equal(t, question.Order, result.Order)
	assert.Equal(t, question.Category, result.Category)
	assert.Equal(t, question.Order, result.Order)
	assert.Equal(t, question.DurationInSeconds, result.DurationInSeconds)
	assert.Equal(t, question.Options, result.Options)
}

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
