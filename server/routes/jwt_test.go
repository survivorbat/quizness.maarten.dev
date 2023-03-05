package routes

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/survivorbat/qq.maarten.dev/server/services"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJwtHandler_JwtGuard_SetsUserOnSuccessful(t *testing.T) {
	t.Parallel()
	// Arrange
	mockJwtService := &MockJwtService{
		validateTokenReturns: &jwt.Token{
			Valid: true,
			Claims: services.QQClaims{
				UserID: "2f80947c-e724-4b38-8c8d-3823864fef58",
			},
		},
	}
	jwtHandler := &JwtHandler{
		JwtService: mockJwtService,
	}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Request, _ = http.NewRequest("", "", nil)
	context.Request.Header = http.Header{"Authorization": []string{"Bearer abc"}}

	// Act
	jwtHandler.JwtGuard()(context)

	// Assert
	assert.Equal(t, "abc", mockJwtService.validateTokenCalledWith)

	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, "2f80947c-e724-4b38-8c8d-3823864fef58", context.GetString("user"))
}

func TestJwtHandler_JwtGuard_ReturnsErrorOnMissingHeader(t *testing.T) {
	t.Parallel()
	// Arrange
	jwtHandler := &JwtHandler{}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Request, _ = http.NewRequest("", "", nil)

	// Act
	jwtHandler.JwtGuard()(context)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, writer.Code)
}

func TestJwtHandler_JwtGuard_ReturnsErrorOnMalformedHeader(t *testing.T) {
	t.Parallel()
	// Arrange
	jwtHandler := &JwtHandler{}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Request, _ = http.NewRequest("", "", nil)
	context.Request.Header = http.Header{"Authorization": []string{"a"}}

	// Act
	jwtHandler.JwtGuard()(context)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, writer.Code)
}

func TestJwtHandler_JwtGuard_ReturnsErrorOnValidationError(t *testing.T) {
	t.Parallel()
	// Arrange
	mockJwtService := &MockJwtService{
		validateTokenReturnsError: assert.AnError,
	}
	jwtHandler := &JwtHandler{JwtService: mockJwtService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Request, _ = http.NewRequest("", "", nil)
	context.Request.Header = http.Header{"Authorization": []string{"Bearer abc"}}

	// Act
	jwtHandler.JwtGuard()(context)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, writer.Code)
}

func TestJwtHandler_JwtGuard_ReturnsErrorOnTokenInvalid(t *testing.T) {
	t.Parallel()
	// Arrange
	mockJwtService := &MockJwtService{
		validateTokenReturns: &jwt.Token{
			Valid: false,
		},
	}
	jwtHandler := &JwtHandler{JwtService: mockJwtService}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Request, _ = http.NewRequest("", "", nil)
	context.Request.Header = http.Header{"Authorization": []string{"Bearer abc"}}

	// Act
	jwtHandler.JwtGuard()(context)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, writer.Code)
}
