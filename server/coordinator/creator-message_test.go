package coordinator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreatorMessage_IsValid_ReturnsValidOnGoodType(t *testing.T) {
	t.Parallel()
	// Arrange
	message := &CreatorMessage{Action: NextQuestionAction}

	// Act
	result := message.IsValid()

	// Assert
	assert.True(t, result)
}
func TestCreatorMessage_IsValid_ReturnsInvalidOnBadType(t *testing.T) {
	t.Parallel()
	// Arrange
	message := &CreatorMessage{Action: "no"}

	// Act
	result := message.IsValid()

	// Assert
	assert.False(t, result)
}
