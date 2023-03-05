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
)

func TestQuizHandler_Get_ReturnsExpectedData(t *testing.T) {
	t.Parallel()
	// Arrange
	quizzes := []*domain.Quiz{
		{BaseObject: domain.BaseObject{ID: uuid.MustParse("2f80947c-e724-4b38-8c8d-3823864fef58")}, CreatorID: uuid.MustParse("d2f584da-d340-459a-a1ce-8652446a86ef")},
		{BaseObject: domain.BaseObject{ID: uuid.MustParse("adeb8482-4eb2-4c2d-8eec-97f705260fa8")}, CreatorID: uuid.MustParse("d2f584da-d340-459a-a1ce-8652446a86ef")},
	}

	mockQuizService := &MockQuizService{getByCreatorReturns: quizzes}
	handler := &QuizHandler{
		QuizService: mockQuizService,
	}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", quizzes[0].CreatorID.String())

	// Act
	handler.Get(context)

	// Assert
	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, quizzes[0].CreatorID, mockQuizService.getByCreatorCalledWith)

	var result []*domain.Quiz
	if err := json.Unmarshal(writer.Body.Bytes(), &result); err != nil {
		t.Fatal(err.Error())
	}

	assert.ElementsMatch(t, quizzes, result)
}

func TestQuizHandler_Get_ReturnsErrorOnFetchError(t *testing.T) {
	t.Parallel()
	// Arrange
	mockQuizService := &MockQuizService{getByCreatorReturnsError: assert.AnError}
	handler := &QuizHandler{
		QuizService: mockQuizService,
	}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58") // Different

	// Act
	handler.Get(context)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, writer.Code)
}
