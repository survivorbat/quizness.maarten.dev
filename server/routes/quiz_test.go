package routes

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/survivorbat/qq.maarten.dev/server/domain"
	"github.com/survivorbat/qq.maarten.dev/server/routes/inputs"
	"gorm.io/gorm"
	"io"
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
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")

	// Act
	handler.Get(context)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, writer.Code)
}

func TestQuizHandler_Post_ReturnsValidationError(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := &QuizHandler{}

	input := &inputs.Quiz{}
	inputJson, _ := json.Marshal(input)

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Request, _ = http.NewRequest(http.MethodPost, "", io.NopCloser(bytes.NewBuffer(inputJson)))

	// Act
	handler.Post(context)

	// Assert
	assert.Equal(t, http.StatusBadRequest, writer.Code)
}

func TestQuizHandler_Post_CallsService(t *testing.T) {
	//t.Parallel() Can't be run in parallel because of the override
	// Arrange
	quizService := &MockQuizService{}
	handler := &QuizHandler{QuizService: quizService}

	input := &inputs.Quiz{
		Name:        "My Awesome Quiz",
		Description: "Best quiz ever",
		MultipleChoiceQuestions: []*inputs.MultipleChoiceQuestion{
			{
				Title:             "What is 2+2",
				Description:       "Simple math",
				DurationInSeconds: 15,
				Category:          "Math",
				Order:             1,
				Options: []*inputs.QuestionOption{
					{TextOption: "20"},
					{TextOption: "15"},
					{TextOption: "4", Answer: true},
					{TextOption: "3"},
				},
			},
		},
	}
	inputJson, _ := json.Marshal(input)

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Request, _ = http.NewRequest(http.MethodPost, "", io.NopCloser(bytes.NewBuffer(inputJson)))

	answerIDs := []uuid.UUID{
		uuid.MustParse("2f80947c-e724-4b38-8c8d-3823864fef58"),
		uuid.MustParse("044bc72e-145a-47fb-969c-16577a08c0e4"),
		uuid.MustParse("76e86314-ce95-4a47-98ce-180dcc724432"),
		uuid.MustParse("1ccca86d-ecb7-4a3e-8097-0e30a3a404e1"),
	}

	var newUuidCalls int
	newUuid = func() uuid.UUID {
		newUuidCalls++

		if len(answerIDs) >= newUuidCalls {
			return answerIDs[newUuidCalls-1]
		}

		return uuid.New()
	}

	// Act
	handler.Post(context)

	// Assert
	assert.Equal(t, http.StatusOK, writer.Code)

	expected := &domain.Quiz{
		Name:        "My Awesome Quiz",
		Description: "Best quiz ever",
		CreatorID:   uuid.MustParse("2f80947c-e724-4b38-8c8d-3823864fef58"),
		MultipleChoiceQuestions: []*domain.MultipleChoiceQuestion{
			{
				BaseQuestion: domain.BaseQuestion{
					Title:             "What is 2+2",
					Description:       "Simple math",
					DurationInSeconds: 15,
					Category:          "Math",
					Order:             1,
				},
				AnswerID: answerIDs[2],
				Options: []*domain.QuestionOption{
					{TextOption: "20", BaseObject: domain.BaseObject{ID: answerIDs[0]}},
					{TextOption: "15", BaseObject: domain.BaseObject{ID: answerIDs[1]}},
					{TextOption: "4", BaseObject: domain.BaseObject{ID: answerIDs[2]}},
					{TextOption: "3", BaseObject: domain.BaseObject{ID: answerIDs[3]}},
				},
			},
		},
	}

	assert.Equal(t, expected, quizService.createCalledWith)
}

func TestQuizHandler_Post_ReturnsAnyErrors(t *testing.T) {
	t.Parallel()
	// Arrange
	quizService := &MockQuizService{createReturns: assert.AnError}
	handler := &QuizHandler{QuizService: quizService}

	input := &inputs.Quiz{
		Name:        "My Awesome Quiz",
		Description: "Best quiz ever",
		MultipleChoiceQuestions: []*inputs.MultipleChoiceQuestion{
			{
				Title:             "What is 2+2",
				Description:       "Simple math",
				DurationInSeconds: 15,
				Category:          "Math",
				Order:             1,
				Options: []*inputs.QuestionOption{
					{TextOption: "20"},
					{TextOption: "15"},
					{TextOption: "4", Answer: true},
					{TextOption: "3"},
				},
			},
		},
	}
	inputJson, _ := json.Marshal(input)

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Request, _ = http.NewRequest(http.MethodPost, "", io.NopCloser(bytes.NewBuffer(inputJson)))

	// Act
	handler.Post(context)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, writer.Code)
}

func TestQuizHandler_Delete_ReturnsParseUUIDError(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := &QuizHandler{}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "cf4851dc-21ec-4eac-8168-4978a40cbc4b")
	context.Params = []gin.Param{{Key: "id", Value: "no"}}

	// Act
	handler.Delete(context)

	// Assert
	assert.Equal(t, http.StatusBadRequest, writer.Code)
}

func TestQuizHandler_Delete_ReturnsFetchByIDError(t *testing.T) {
	t.Parallel()
	// Arrange
	service := &MockQuizService{
		getByIdReturnsError: gorm.ErrRecordNotFound,
	}
	handler := &QuizHandler{QuizService: service}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Params = []gin.Param{{Key: "id", Value: "e54f8551-4352-48e9-b6f4-2c293b31de9c"}}

	// Act
	handler.Delete(context)

	// Assert
	assert.Equal(t, http.StatusNotFound, writer.Code)
}

func TestQuizHandler_Delete_ReturnsErrorOnNotMyQuiz(t *testing.T) {
	t.Parallel()
	// Arrange
	service := &MockQuizService{
		getByIdReturns: &domain.Quiz{
			CreatorID: uuid.MustParse("67389c72-b059-4680-9e76-bad17b6d40c5"),
		},
	}
	handler := &QuizHandler{QuizService: service}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Params = []gin.Param{{Key: "id", Value: "e54f8551-4352-48e9-b6f4-2c293b31de9c"}}

	// Act
	handler.Delete(context)

	// Assert
	assert.Equal(t, http.StatusForbidden, writer.Code)
}

func TestQuizHandler_Delete_ReturnsErrorOnDeleteFailure(t *testing.T) {
	t.Parallel()
	// Arrange
	service := &MockQuizService{
		getByIdReturns: &domain.Quiz{
			BaseObject: domain.BaseObject{ID: uuid.MustParse("e54f8551-4352-48e9-b6f4-2c293b31de9c")},
			CreatorID:  uuid.MustParse("2f80947c-e724-4b38-8c8d-3823864fef58"),
		},
		deleteReturns: assert.AnError,
	}
	handler := &QuizHandler{QuizService: service}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Params = []gin.Param{{Key: "id", Value: "e54f8551-4352-48e9-b6f4-2c293b31de9c"}}

	// Act
	handler.Delete(context)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, writer.Code)
}

func TestQuizHandler_Delete_CallsDelete(t *testing.T) {
	t.Parallel()
	// Arrange
	service := &MockQuizService{
		getByIdReturns: &domain.Quiz{
			BaseObject: domain.BaseObject{ID: uuid.MustParse("e54f8551-4352-48e9-b6f4-2c293b31de9c")},
			CreatorID:  uuid.MustParse("2f80947c-e724-4b38-8c8d-3823864fef58"),
		},
	}
	handler := &QuizHandler{QuizService: service}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58")
	context.Params = []gin.Param{{Key: "id", Value: "e54f8551-4352-48e9-b6f4-2c293b31de9c"}}

	// Act
	handler.Delete(context)

	// Assert
	assert.Equal(t, http.StatusNoContent, writer.Code)
	assert.Equal(t, service.getByIdReturns.ID, service.deleteCalledWith)
}
