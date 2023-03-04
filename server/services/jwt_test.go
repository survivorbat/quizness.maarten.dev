package services

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJwtService_GenerateToken_ReturnsExpectedToken(t *testing.T) {
	t.Parallel()
	// Arrange
	service := &JwtService{SecretKey: "abc", Issuer: "My Company"}
	userID := uuid.MustParse("6aacfb41-e478-46ec-857e-11221f2a97fc")

	// Act
	token, err := service.GenerateToken(userID)

	// Assert
	assert.NoError(t, err)
	assert.Contains(t, token, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9")
}

func TestJwtService_ValidateToken_ReturnsValidOnGoodToken(t *testing.T) {
	t.Parallel()
	// Arrange
	service := &JwtService{SecretKey: "abc", Issuer: "My Company"}
	userID := uuid.MustParse("6aacfb41-e478-46ec-857e-11221f2a97fc")

	token, _ := service.GenerateToken(userID)

	// Act
	result, err := service.ValidateToken(token)

	// Assert
	assert.NoError(t, err)
	assert.True(t, result.Valid)
}

func TestJwtService_ValidateToken_ReturnsInvalidOnBadToken(t *testing.T) {
	t.Parallel()
	// Arrange
	service := &JwtService{SecretKey: "abc", Issuer: "My Company"}

	// Act
	result, err := service.ValidateToken("abc")

	// Assert
	assert.ErrorContains(t, err, "token contains an invalid number of segments")
	assert.Empty(t, result)
}
