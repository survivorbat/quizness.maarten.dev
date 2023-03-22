package domain

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestCreator_EnsureHasNickname_SetsARandomName(t *testing.T) {
	t.Parallel()
	// Arrange
	creator := &Creator{}

	// Act
	creator.GenerateNickname()

	// Assert
	split := strings.Split(creator.Nickname, " ")
	assert.Contains(t, namePrefixes, split[0])
	assert.Contains(t, nameSuffixes, split[1])
}

func TestCreator_GenerateColors_SetsARandomColor(t *testing.T) {
	t.Parallel()
	// Arrange
	creator := &Creator{}

	// Act
	creator.GenerateColors()

	// Assert
	assert.Len(t, creator.Color, 7)
	assert.Len(t, creator.BackgroundColor, 7)
}
