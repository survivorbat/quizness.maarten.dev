package outputs

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/survivorbat/qq.maarten.dev/server/domain"
	"testing"
)

func TestNewPublicGame_ReturnsExpectedValue(t *testing.T) {
	t.Parallel()
	// Arrange
	game := &domain.Game{BaseObject: domain.BaseObject{ID: uuid.MustParse("942ee70d-1d18-4b1d-8abb-fdb696d2da0b")}}

	// Act
	result := NewPublicGame(game)

	// Assert
	assert.Equal(t, game.ID, result.ID)
}
