package routes

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/survivorbat/qq.maarten.dev/server/domain"
	"github.com/survivorbat/qq.maarten.dev/server/routes/inputs"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAnswerHandler_Patch_ReturnsErrorOnInvalidGameUUID(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := &AnswerHandler{}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Request, _ = http.NewRequest(http.MethodPatch, "", nil)
	context.Params = []gin.Param{{Key: "id", Value: "no"}}

	// Act
	handler.Patch(context)

	// Assert
	assert.Equal(t, http.StatusBadRequest, writer.Code)
}

func TestAnswerHandler_Patch_ReturnsErrorOnInvalidQuestionUUID(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := &AnswerHandler{}

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

func TestAnswerHandler_Patch_ReturnsErrorOnInvalidPlayerUUID(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := &AnswerHandler{}

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

func TestAnswerHandler_Patch_ReturnsErrorOnGameNotFound(t *testing.T) {
	t.Parallel()
	// Arrange
	gameService := &MockGameService{getByIdReturnsError: assert.AnError}
	handler := &AnswerHandler{GameService: gameService}

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

func TestAnswerHandler_Patch_ReturnsValidationError(t *testing.T) {
	t.Parallel()
	// Arrange
	gameService := &MockGameService{getByIdReturns: &domain.Game{}}
	handler := &AnswerHandler{GameService: gameService}

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

func TestAnswerHandler_Patch_ReturnsGenericError(t *testing.T) {
	t.Parallel()
	// Arrange
	gameService := &MockGameService{getByIdReturns: &domain.Game{}, answerReturns: assert.AnError}
	handler := &AnswerHandler{GameService: gameService}

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

func TestAnswerHandler_Patch_ReturnsSuccess(t *testing.T) {
	t.Parallel()
	// Arrange
	gameService := &MockGameService{getByIdReturns: &domain.Game{}}
	handler := &AnswerHandler{GameService: gameService}

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
