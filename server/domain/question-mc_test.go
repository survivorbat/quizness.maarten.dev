package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMultipleChoiceQuestion_GetType_ReturnsExpected(t *testing.T) {
	t.Parallel()
	// Arrange
	mc := new(MultipleChoiceQuestion)

	// Act
	result := mc.GetType()

	// Assert
	assert.Equal(t, TypeMultipleChoice, result)
}
