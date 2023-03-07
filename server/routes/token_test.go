package routes

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTokenHandler_JwtGuard_SetsUserOnSuccessful(t *testing.T) {
	t.Parallel()
	// Arrange
	mockJwtService := &MockJwtService{
		validateTokenReturns: &jwt.Token{
			Valid: true,
			Claims: jwt.MapClaims{
				"userID": "2f80947c-e724-4b38-8c8d-3823864fef58",
			},
		},
	}
	tokenHandler := &TokenHandler{
		JwtService: mockJwtService,
	}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Request, _ = http.NewRequest("", "", nil)
	context.Request.Header = http.Header{"Authorization": []string{"Bearer abc"}}

	// Act
	tokenHandler.JwtGuard()(context)

	// Assert
	assert.Equal(t, "abc", mockJwtService.validateTokenCalledWith)

	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, "2f80947c-e724-4b38-8c8d-3823864fef58", context.GetString("user"))
}

func TestTokenHandler_JwtGuard_ReturnsErrorOnMissingHeader(t *testing.T) {
	t.Parallel()
	// Arrange
	tokenHandler := &TokenHandler{}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Request, _ = http.NewRequest("", "", nil)

	// Act
	tokenHandler.JwtGuard()(context)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, writer.Code)
}

func TestTokenHandler_JwtGuard_ReturnsErrorOnMalformedHeader(t *testing.T) {
	t.Parallel()
	// Arrange
	tokenHandler := &TokenHandler{}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Request, _ = http.NewRequest("", "", nil)
	context.Request.Header = http.Header{"Authorization": []string{"a"}}

	// Act
	tokenHandler.JwtGuard()(context)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, writer.Code)
}

func TestTokenHandler_JwtGuard_ReturnsErrorOnValidationError(t *testing.T) {
	t.Parallel()
	// Arrange
	mockJwtService := &MockJwtService{
		validateTokenReturnsError: assert.AnError,
	}
	tokenHandler := &TokenHandler{JwtService: mockJwtService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Request, _ = http.NewRequest("", "", nil)
	context.Request.Header = http.Header{"Authorization": []string{"Bearer abc"}}

	// Act
	tokenHandler.JwtGuard()(context)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, writer.Code)
}

func TestTokenHandler_JwtGuard_ReturnsErrorOnTokenInvalid(t *testing.T) {
	t.Parallel()
	// Arrange
	mockJwtService := &MockJwtService{
		validateTokenReturns: &jwt.Token{
			Valid: false,
		},
	}
	tokenHandler := &TokenHandler{JwtService: mockJwtService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Request, _ = http.NewRequest("", "", nil)
	context.Request.Header = http.Header{"Authorization": []string{"Bearer abc"}}

	// Act
	tokenHandler.JwtGuard()(context)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, writer.Code)
}
