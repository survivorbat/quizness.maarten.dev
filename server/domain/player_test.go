package domain

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestPlayer_EnsureHasNickname_SetsARandomName(t *testing.T) {
	t.Parallel()
	// Arrange
	player := &Player{}

	// Act
	player.GenerateNickname()

	// Assert
	split := strings.Split(player.Nickname, " ")
	assert.Contains(t, namePrefixes, split[0])
	assert.Contains(t, nameSuffixes, split[1])
}
