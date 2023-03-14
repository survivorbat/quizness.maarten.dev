package domain

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestPlayers_Contains_ReturnsTrueOnContains(t *testing.T) {
	t.Parallel()
	// Arrange
	players := Players{
		{BaseObject: BaseObject{ID: uuid.MustParse("38927e3b-e2ef-4207-9a9b-bd079928b1f0")}},
		{BaseObject: BaseObject{ID: uuid.MustParse("32a7b8d7-3c9b-49c9-bdd1-9c5415a059fc")}},
		{BaseObject: BaseObject{ID: uuid.MustParse("1b00cbc4-e647-4868-8953-70fc559b8ff0")}},
	}

	// Act
	result := players.Contains(players[2].ID)

	// Assert
	assert.True(t, result)
}

func TestPlayers_Contains_ReturnsFalseOnNotContains(t *testing.T) {
	t.Parallel()
	// Arrange
	players := Players{
		{BaseObject: BaseObject{ID: uuid.MustParse("20ac1eb2-5424-4aca-9018-1eee6a2a5510")}},
		{BaseObject: BaseObject{ID: uuid.MustParse("e2c0e26c-91eb-4946-8ea4-5304036df91c")}},
		{BaseObject: BaseObject{ID: uuid.MustParse("d3123f3a-b4e2-45c1-8e15-8696fcebc7c4")}},
	}

	// Act
	result := players.Contains(uuid.MustParse("888d6988-d39d-4f75-9987-499f5cf12375"))

	// Assert
	assert.False(t, result)
}

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
