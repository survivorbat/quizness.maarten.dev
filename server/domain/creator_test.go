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
