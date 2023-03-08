package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGame_PlayerJoin_AddsPlayer(t *testing.T) {
	t.Parallel()
	// Arrange
	player := &Player{Nickname: "abc"}
	game := Game{}

	// Act
	err := game.PlayerJoin(player)

	// Assert
	assert.NoError(t, err)
}

func TestGame_PlayerJoin_ReturnsErrorOnAlreadyPresent(t *testing.T) {
	t.Parallel()
	// Arrange
	player := &Player{Nickname: "abc"}
	game := Game{Players: []*Player{{Nickname: "abc"}}}

	// Act
	err := game.PlayerJoin(player)

	// Assert
	assert.ErrorContains(t, err, "player is already in this game")
}

func TestGame_PlayerLeave_AddsPlayer(t *testing.T) {
	t.Parallel()
	// Arrange
	player := &Player{Nickname: "abc"}
	game := Game{Players: []*Player{{Nickname: "abc"}}}

	// Act
	err := game.PlayerLeave(player)

	// Assert
	assert.NoError(t, err)
}

func TestGame_PlayerLeave_ReturnsErrorOnAlreadyPresent(t *testing.T) {
	t.Parallel()
	// Arrange
	player := &Player{Nickname: "abc"}
	game := Game{Players: []*Player{}}

	// Act
	err := game.PlayerLeave(player)

	// Assert
	assert.ErrorContains(t, err, "player is not in this game")
}
