package coordinator

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/survivorbat/qq.maarten.dev/server/inputs"
	"testing"
)

func TestPlayerMessage_IsValid_ReturnsExpectedValues(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		message  *PlayerMessage
		expected bool
	}{
		"invalid action": {
			message: &PlayerMessage{Action: "no"},
		},
		"invalid answer": {
			message: &PlayerMessage{
				Action: AnswerAction,
				Answer: &inputs.Answer{},
			},
		},
		"valid answer": {
			message: &PlayerMessage{
				Action: AnswerAction,
				Answer: &inputs.Answer{OptionID: uuid.MustParse("5b8c33ef-75cf-4508-9ab7-952dfd1ed240")},
			},
			expected: true,
		},
	}

	for name, testData := range tests {
		testData := testData
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			// Act
			result := testData.message.IsValid()

			// Assert
			assert.Equal(t, testData.expected, result)
		})
	}
}

func TestPlayerMessage_Parse_ParsesCorrectly(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		message  *PlayerMessage
		content  any
		expected *PlayerMessage
	}{
		"answer": {
			message: &PlayerMessage{
				Action: AnswerAction,
			},
			content: &inputs.Answer{
				OptionID: uuid.MustParse("5b8c33ef-75cf-4508-9ab7-952dfd1ed240"),
			},
			expected: &PlayerMessage{
				Action: AnswerAction,
				Answer: &inputs.Answer{
					OptionID: uuid.MustParse("5b8c33ef-75cf-4508-9ab7-952dfd1ed240"),
				},
			},
		},
	}

	for name, testData := range tests {
		testData := testData
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			testData.message.Content, _ = json.Marshal(testData.content)
			testData.expected.Content = testData.message.Content

			// Act
			err := testData.message.Parse()

			// Assert
			assert.NoError(t, err)

			assert.Equal(t, testData.expected, testData.message)
		})
	}
}

func TestPlayerMessage_Parse_ReturnsErrorOnInvalidAction(t *testing.T) {
	t.Parallel()
	// Arrange
	message := &PlayerMessage{Action: "no"}

	// Act
	err := message.Parse()

	// Assert
	assert.ErrorContains(t, err, "unknown type")
}
