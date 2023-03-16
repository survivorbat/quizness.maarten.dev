package routes

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/survivorbat/qq.maarten.dev/server/domain"
	"github.com/survivorbat/qq.maarten.dev/server/inputs"
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
	context.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	context.Params = []gin.Param{
		{Key: "code", Value: "abc"},
	}

	// Act
	handler.GetByCode(context)

	// Assert
	assert.Equal(t, http.StatusNotFound, writer.Code)
}

func TestPublicGameHandler_GetByCode_ReturnsGameID(t *testing.T) {
	t.Parallel()
	// Arrange
	game := &domain.Game{BaseObject: domain.BaseObject{ID: uuid.MustParse("788f12a9-51e8-4c87-9b0c-06bcc9f0691b")}}
	gameService := &MockGameService{getByCodeReturns: game}
	handler := &PublicGameHandler{GameService: gameService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	context.Params = []gin.Param{
		{Key: "code", Value: "abc"},
	}

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

func TestPublicGameHandler_Patch_ReturnsErrorOnInvalidGameUUID(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := &PublicGameHandler{}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Request, _ = http.NewRequest(http.MethodPatch, "", nil)
	context.Params = []gin.Param{{Key: "id", Value: "no"}}

	// Act
	handler.Patch(context)

	// Assert
	assert.Equal(t, http.StatusBadRequest, writer.Code)
}

func TestPublicGameHandler_Patch_ReturnsErrorOnInvalidQuestionUUID(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := &PublicGameHandler{}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Request, _ = http.NewRequest(http.MethodPatch, "", nil)
	context.Params = []gin.Param{
		{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"},
	}

	// Act
	handler.Patch(context)

	// Assert
	assert.Equal(t, http.StatusBadRequest, writer.Code)
}

func TestPublicGameHandler_Patch_ReturnsErrorOnInvalidPlayerUUID(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := &PublicGameHandler{}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Request, _ = http.NewRequest(http.MethodPatch, "", nil)
	context.Params = []gin.Param{
		{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"},
		{Key: "question", Value: "1e8df2b3-2cf7-42d7-92cf-9826a5697c69"},
	}

	// Act
	handler.Patch(context)

	// Assert
	assert.Equal(t, http.StatusBadRequest, writer.Code)
}

func TestPublicGameHandler_Patch_ReturnsErrorOnGameNotFound(t *testing.T) {
	t.Parallel()
	// Arrange
	gameService := &MockGameService{getByIdReturnsError: assert.AnError}
	handler := &PublicGameHandler{GameService: gameService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Request, _ = http.NewRequest(http.MethodPatch, "", nil)
	context.Params = []gin.Param{
		{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"},
		{Key: "question", Value: "1e8df2b3-2cf7-42d7-92cf-9826a5697c69"},
		{Key: "player", Value: "707dd608-7e76-4993-a573-4ad6f809cb7c"},
	}

	// Act
	handler.Patch(context)

	// Assert
	assert.Equal(t, http.StatusNotFound, writer.Code)
}

func TestPublicGameHandler_Patch_ReturnsValidationError(t *testing.T) {
	t.Parallel()
	// Arrange
	gameService := &MockGameService{getByIdReturns: &domain.Game{}}
	handler := &PublicGameHandler{GameService: gameService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Request, _ = http.NewRequest(http.MethodPatch, "", nil)
	context.Params = []gin.Param{
		{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"},
		{Key: "question", Value: "1e8df2b3-2cf7-42d7-92cf-9826a5697c69"},
		{Key: "player", Value: "707dd608-7e76-4993-a573-4ad6f809cb7c"},
	}

	// Act
	handler.Patch(context)

	// Assert
	assert.Equal(t, http.StatusBadRequest, writer.Code)
}

func TestPublicGameHandler_Patch_ReturnsGenericError(t *testing.T) {
	t.Parallel()
	// Arrange
	gameService := &MockGameService{getByIdReturns: &domain.Game{}, answerReturns: assert.AnError}
	handler := &PublicGameHandler{GameService: gameService}

	input := &inputs.Answer{OptionID: uuid.MustParse("30f9b929-e202-4305-ac30-114875c19cc3")}
	inputJson, _ := json.Marshal(input)

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Request, _ = http.NewRequest(http.MethodPatch, "", io.NopCloser(bytes.NewBuffer(inputJson)))
	context.Params = []gin.Param{
		{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"},
		{Key: "question", Value: "1e8df2b3-2cf7-42d7-92cf-9826a5697c69"},
		{Key: "player", Value: "707dd608-7e76-4993-a573-4ad6f809cb7c"},
	}

	// Act
	handler.Patch(context)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, writer.Code)
}

func TestPublicGameHandler_Patch_ReturnsSuccess(t *testing.T) {
	t.Parallel()
	// Arrange
	gameService := &MockGameService{getByIdReturns: &domain.Game{}}
	handler := &PublicGameHandler{GameService: gameService}

	input := &inputs.Answer{OptionID: uuid.MustParse("30f9b929-e202-4305-ac30-114875c19cc3")}
	inputJson, _ := json.Marshal(input)

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Request, _ = http.NewRequest(http.MethodPatch, "", io.NopCloser(bytes.NewBuffer(inputJson)))
	context.Params = []gin.Param{
		{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"},
		{Key: "question", Value: "1e8df2b3-2cf7-42d7-92cf-9826a5697c69"},
		{Key: "player", Value: "707dd608-7e76-4993-a573-4ad6f809cb7c"},
	}

	// Act
	handler.Patch(context)

	// Assert
	assert.Equal(t, http.StatusOK, writer.Code)
}
