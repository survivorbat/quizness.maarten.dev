package coordinator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlayerAction_IsValid_ReturnsExpectedValue(t *testing.T) {
	t.Parallel()
	tests := map[PlayerAction]bool{
		AnswerAction: true,
		"no":         false,
	}

	for action, expected := range tests {
		expected := expected
		action := action
		t.Run(string(action), func(t *testing.T) {
			t.Parallel()
			// Act
			result := action.IsValid()

			// Assert
			assert.Equal(t, expected, result)
		})
	}
}

func TestCreatorAction_IsValid_ReturnsExpectedValue(t *testing.T) {
	t.Parallel()
	tests := map[CreatorAction]bool{
		NextQuestionAction: true,
		FinishGameAction:   true,
		"no":               false,
	}

	for action, expected := range tests {
		expected := expected
		action := action
		t.Run(string(action), func(t *testing.T) {
			t.Parallel()
			// Act
			result := action.IsValid()

			// Assert
			assert.Equal(t, expected, result)
		})
	}
}
