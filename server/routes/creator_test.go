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

func TestCreatorHandler_GetWithID_ReturnsExpectedData(t *testing.T) {
	t.Parallel()
	// Arrange
	creator := &domain.Creator{
		BaseObject: domain.BaseObject{ID: uuid.MustParse("2f80947c-e724-4b38-8c8d-3823864fef58")},
		Nickname:   "abc",
		AuthID:     "def",
	}

	mockCreatorService := &MockCreatorService{getByIDReturns: creator}
	handler := &CreatorHandler{
		CreatorService: mockCreatorService,
	}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Set("user", creator.ID.String())

	// Act
	handler.GetWithID(context)

	// Assert
	assert.Equal(t, http.StatusOK, writer.Code)

	assert.Equal(t, creator.ID, mockCreatorService.getByIDCalledWith)

	var result *domain.Creator
	if err := json.Unmarshal(writer.Body.Bytes(), &result); err != nil {
		t.Fatal(err.Error())
	}

	assert.Equal(t, creator.ID, result.ID)
	assert.Equal(t, creator.Nickname, result.Nickname)
	assert.Empty(t, result.AuthID)
}

func TestCreatorHandler_GetWithID_ReturnsErrorOnFetchError(t *testing.T) {
	t.Parallel()
	// Arrange
	creator := &domain.Creator{
		BaseObject: domain.BaseObject{ID: uuid.MustParse("2f80947c-e724-4b38-8c8d-3823864fef58")},
		Nickname:   "abc",
		AuthID:     "def",
	}

	mockCreatorService := &MockCreatorService{getByIDReturnsError: assert.AnError}
	handler := &CreatorHandler{
		CreatorService: mockCreatorService,
	}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Params = gin.Params{{Key: "id", Value: creator.ID.String()}}
	context.Set("user", "2f80947c-e724-4b38-8c8d-3823864fef58") // Different

	// Act
	handler.GetWithID(context)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, writer.Code)
}
