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
	"time"
)

func TestGameHandler_Get_ReturnsErrorOnInvalidUUID(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := &GameHandler{}

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

func TestGameHandler_Get_ReturnsErrorOnQuizNotFound(t *testing.T) {
	t.Parallel()
	// Arrange
	quizService := &MockQuizService{getByIdReturnsError: assert.AnError}
	handler := &GameHandler{QuizService: quizService}

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

func TestGameHandler_Get_ReturnsErrorOnNotMyQuiz(t *testing.T) {
	t.Parallel()
	// Arrange
	quizService := &MockQuizService{getByIdReturns: &domain.Quiz{
		CreatorID: uuid.MustParse("8fdc3e5a-b0a8-4103-af3b-c2f20d91889b"),
	}}
	handler := &GameHandler{QuizService: quizService}

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

func TestGameHandler_Get_ReturnsErrorOnFailedGet(t *testing.T) {
	t.Parallel()
	// Arrange
	quizService := &MockQuizService{getByIdReturns: &domain.Quiz{
		CreatorID: uuid.MustParse("2f80947c-e724-4b38-8c8d-3823864fef58"),
	}}
	gameService := &MockGameService{getByQuizReturnsError: assert.AnError}
	handler := &GameHandler{QuizService: quizService, GameService: gameService}

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

func TestGameHandler_Get_ReturnsData(t *testing.T) {
	t.Parallel()
	// Arrange
	quizService := &MockQuizService{getByIdReturns: &domain.Quiz{
		CreatorID: uuid.MustParse("2f80947c-e724-4b38-8c8d-3823864fef58"),
	}}
	gameService := &MockGameService{getByQuizReturns: []*domain.Game{{Code: "A"}, {Code: "B"}}}
	handler := &GameHandler{QuizService: quizService, GameService: gameService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	context.Params = []gin.Param{{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"}}

	// Act
	handler.Get(context)

	// Assert
	assert.Equal(t, http.StatusOK, writer.Code)

	var result []*domain.Game
	if err := json.Unmarshal(writer.Body.Bytes(), &result); err != nil {
		t.Fatal(err.Error())
	}

	assert.ElementsMatch(t, gameService.getByQuizReturns, result)
}

func TestGameHandler_Post_ReturnsErrorOnInvalidUUID(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := &GameHandler{}

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

func TestGameHandler_Post_ReturnsErrorOnQuizNotFound(t *testing.T) {
	t.Parallel()
	// Arrange
	quizService := &MockQuizService{getByIdReturnsError: assert.AnError}
	handler := &GameHandler{QuizService: quizService}

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

func TestGameHandler_Post_ReturnsErrorOnNotMyQuiz(t *testing.T) {
	t.Parallel()
	// Arrange
	quizService := &MockQuizService{getByIdReturns: &domain.Quiz{
		CreatorID: uuid.MustParse("8fdc3e5a-b0a8-4103-af3b-c2f20d91889b"),
	}}
	handler := &GameHandler{QuizService: quizService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Request, _ = http.NewRequest(http.MethodPost, "", nil)
	context.Params = []gin.Param{{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"}}

	// Act
	handler.Post(context)

	// Assert
	assert.Equal(t, http.StatusForbidden, writer.Code)
}

func TestGameHandler_Post_ReturnsErrorOnValidationErrors(t *testing.T) {
	t.Parallel()
	// Arrange
	quizService := &MockQuizService{getByIdReturns: &domain.Quiz{
		CreatorID: uuid.MustParse("2f80947c-e724-4b38-8c8d-3823864fef58"),
	}}
	handler := &GameHandler{QuizService: quizService}

	input := &inputs.Game{PlayerLimit: 0}
	inputJson, _ := json.Marshal(input)

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Request, _ = http.NewRequest(http.MethodPost, "", io.NopCloser(bytes.NewBuffer(inputJson)))
	context.Params = []gin.Param{{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"}}

	// Act
	handler.Post(context)

	// Assert
	assert.Equal(t, http.StatusBadRequest, writer.Code)
}
func TestGameHandler_Post_ReturnsGenericErrorOnCreate(t *testing.T) {
	t.Parallel()
	// Arrange
	quizService := &MockQuizService{getByIdReturns: &domain.Quiz{
		CreatorID: uuid.MustParse("2f80947c-e724-4b38-8c8d-3823864fef58"),
	}}
	gameService := &MockGameService{createReturns: assert.AnError}
	handler := &GameHandler{QuizService: quizService, GameService: gameService}

	input := &inputs.Game{PlayerLimit: 2}
	inputJson, _ := json.Marshal(input)

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Request, _ = http.NewRequest(http.MethodPost, "", io.NopCloser(bytes.NewBuffer(inputJson)))
	context.Params = []gin.Param{{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"}}

	// Act
	handler.Post(context)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, writer.Code)
}

func TestGameHandler_Post_ReturnsData(t *testing.T) {
	t.Parallel()
	// Arrange
	quizService := &MockQuizService{getByIdReturns: &domain.Quiz{
		CreatorID: uuid.MustParse("2f80947c-e724-4b38-8c8d-3823864fef58"),
	}}
	gameService := &MockGameService{}
	handler := &GameHandler{QuizService: quizService, GameService: gameService}

	input := &inputs.Game{PlayerLimit: 2}
	inputJson, _ := json.Marshal(input)

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Request, _ = http.NewRequest(http.MethodPost, "", io.NopCloser(bytes.NewBuffer(inputJson)))
	context.Params = []gin.Param{{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"}}

	// Act
	handler.Post(context)

	// Assert
	assert.Equal(t, http.StatusOK, writer.Code)

	var result *domain.Game
	if err := json.Unmarshal(writer.Body.Bytes(), &result); err != nil {
		t.Fatal(err.Error())
	}

	assert.Equal(t, input.PlayerLimit, result.PlayerLimit)
}

func TestGameHandler_Patch_ReturnsErrorOnInvalidUUID(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := &GameHandler{}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Request, _ = http.NewRequest(http.MethodPatch, "", nil)
	context.Params = []gin.Param{{Key: "id", Value: "no"}}

	// Act
	handler.Patch(context)

	// Assert
	assert.Equal(t, http.StatusBadRequest, writer.Code)
}

func TestGameHandler_Patch_ReturnsErrorOnGameNotFound(t *testing.T) {
	t.Parallel()
	// Arrange
	gameService := &MockGameService{getByIdReturnsError: assert.AnError}
	handler := &GameHandler{GameService: gameService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Request, _ = http.NewRequest(http.MethodPatch, "", nil)
	context.Params = []gin.Param{{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"}}

	// Act
	handler.Patch(context)

	// Assert
	assert.Equal(t, http.StatusNotFound, writer.Code)
}

func TestGameHandler_Patch_ReturnsErrorOnNotMyQuiz(t *testing.T) {
	t.Parallel()
	// Arrange
	gameService := &MockGameService{
		getByIdReturns: &domain.Game{Quiz: &domain.Quiz{CreatorID: uuid.MustParse("76afc275-5454-4359-a52b-02693a9c48ba")}},
	}
	handler := &GameHandler{GameService: gameService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Request, _ = http.NewRequest(http.MethodPatch, "", nil)
	context.Params = []gin.Param{{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"}}

	// Act
	handler.Patch(context)

	// Assert
	assert.Equal(t, http.StatusForbidden, writer.Code)
}

func TestGameHandler_Patch_ReturnsErrorInvalidAction(t *testing.T) {
	t.Parallel()
	// Arrange
	gameService := &MockGameService{
		getByIdReturns: &domain.Game{Quiz: &domain.Quiz{CreatorID: uuid.MustParse("2f80947c-e724-4b38-8c8d-3823864fef58")}},
	}
	handler := &GameHandler{GameService: gameService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Request, _ = http.NewRequest(http.MethodPatch, "https://test.com?action=unknown", nil)
	context.Params = []gin.Param{{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"}}

	// Act
	handler.Patch(context)

	// Assert
	assert.Equal(t, http.StatusBadRequest, writer.Code)
}

func TestGameHandler_Patch_ReturnsOnNextError(t *testing.T) {
	t.Parallel()
	// Arrange
	gameService := &MockGameService{
		getByIdReturns: &domain.Game{Quiz: &domain.Quiz{CreatorID: uuid.MustParse("2f80947c-e724-4b38-8c8d-3823864fef58")}},
		nextReturns:    assert.AnError,
	}
	handler := &GameHandler{GameService: gameService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Request, _ = http.NewRequest(http.MethodPatch, "https://test.com?action=next", nil)
	context.Params = []gin.Param{{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"}}

	// Act
	handler.Patch(context)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, writer.Code)
}

func TestGameHandler_Patch_StartReturnsErrorOnGameAlreadyInProgress(t *testing.T) {
	t.Parallel()
	// Arrange
	gameService := &MockGameService{
		getByIdReturns: &domain.Game{
			Quiz: &domain.Quiz{
				CreatorID: uuid.MustParse("2f80947c-e724-4b38-8c8d-3823864fef58"),
				Games:     []*domain.Game{{StartTime: time.Now()}},
			},
		},
	}
	handler := &GameHandler{GameService: gameService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Request, _ = http.NewRequest(http.MethodPatch, "https://test.com?action=start", nil)
	context.Params = []gin.Param{{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"}}

	// Act
	handler.Patch(context)

	// Assert
	assert.Equal(t, http.StatusConflict, writer.Code)
}

func TestGameHandler_Patch_ReturnsStartError(t *testing.T) {
	t.Parallel()
	// Arrange
	gameService := &MockGameService{
		getByIdReturns: &domain.Game{
			Quiz: &domain.Quiz{
				CreatorID: uuid.MustParse("2f80947c-e724-4b38-8c8d-3823864fef58"),
			},
		},
		startReturns: assert.AnError,
	}
	handler := &GameHandler{GameService: gameService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Request, _ = http.NewRequest(http.MethodPatch, "https://test.com?action=start", nil)
	context.Params = []gin.Param{{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"}}

	// Act
	handler.Patch(context)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, writer.Code)
}

func TestGameHandler_Patch_ReturnsFinishError(t *testing.T) {
	t.Parallel()
	// Arrange
	gameService := &MockGameService{
		getByIdReturns: &domain.Game{
			Quiz: &domain.Quiz{
				CreatorID: uuid.MustParse("2f80947c-e724-4b38-8c8d-3823864fef58"),
			},
		},
		finishReturns: assert.AnError,
	}
	handler := &GameHandler{GameService: gameService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Request, _ = http.NewRequest(http.MethodPatch, "https://test.com?action=finish", nil)
	context.Params = []gin.Param{{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"}}

	// Act
	handler.Patch(context)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, writer.Code)
}

func TestGameHandler_Patch_ReturnsResult(t *testing.T) {
	t.Parallel()
	// Arrange
	gameService := &MockGameService{
		getByIdReturns: &domain.Game{
			Quiz: &domain.Quiz{
				CreatorID: uuid.MustParse("2f80947c-e724-4b38-8c8d-3823864fef58"),
			},
		},
	}
	handler := &GameHandler{GameService: gameService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Request, _ = http.NewRequest(http.MethodPatch, "https://test.com?action=start", nil)
	context.Params = []gin.Param{{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"}}

	// Act
	handler.Patch(context)

	// Assert
	assert.Equal(t, http.StatusOK, writer.Code)

	var result *domain.Game
	if err := json.Unmarshal(writer.Body.Bytes(), &result); err != nil {
		t.Fatal(err.Error())
	}

	assert.Equal(t, gameService.getByIdReturns.ID, result.ID)
}

func TestGameHandler_Delete_ReturnsErrorOnInvalidUUID(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := &GameHandler{}

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

func TestGameHandler_Delete_ReturnsErrorOnGameNotFound(t *testing.T) {
	t.Parallel()
	// Arrange
	gameService := &MockGameService{getByIdReturnsError: assert.AnError}
	handler := &GameHandler{GameService: gameService}

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

func TestGameHandler_Delete_ReturnsErrorOnNotMyQuiz(t *testing.T) {
	t.Parallel()
	// Arrange
	gameService := &MockGameService{
		getByIdReturns: &domain.Game{Quiz: &domain.Quiz{CreatorID: uuid.MustParse("76afc275-5454-4359-a52b-02693a9c48ba")}},
	}
	handler := &GameHandler{GameService: gameService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Request, _ = http.NewRequest(http.MethodDelete, "", nil)
	context.Params = []gin.Param{{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"}}

	// Act
	handler.Delete(context)

	// Assert
	assert.Equal(t, http.StatusForbidden, writer.Code)
}

func TestGameHandler_Delete_ReturnsGenericDeleteError(t *testing.T) {
	t.Parallel()
	// Arrange
	gameService := &MockGameService{
		getByIdReturns: &domain.Game{Quiz: &domain.Quiz{CreatorID: uuid.MustParse("2f80947c-e724-4b38-8c8d-3823864fef58")}},
		deleteReturns:  assert.AnError,
	}
	handler := &GameHandler{GameService: gameService}

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

func TestGameHandler_Delete_ReturnsResult(t *testing.T) {
	t.Parallel()
	// Arrange
	gameService := &MockGameService{
		getByIdReturns: &domain.Game{
			Quiz: &domain.Quiz{
				CreatorID: uuid.MustParse("2f80947c-e724-4b38-8c8d-3823864fef58"),
			},
		},
	}
	handler := &GameHandler{GameService: gameService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Request, _ = http.NewRequest(http.MethodDelete, "https://test.com?action=start", nil)
	context.Params = []gin.Param{{Key: "id", Value: "788f12a9-51e8-4c87-9b0c-06bcc9f0691b"}}

	// Act
	handler.Delete(context)

	// Assert
	assert.Equal(t, http.StatusOK, writer.Code)

	var result *domain.Game
	if err := json.Unmarshal(writer.Body.Bytes(), &result); err != nil {
		t.Fatal(err.Error())
	}

	assert.Equal(t, gameService.getByIdReturns.ID, result.ID)
}
