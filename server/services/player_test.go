package services

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/survivorbat/qq.maarten.dev/server/domain"
	"testing"
)

func TestDBPlayerService_GetByID_ReturnsExpected(t *testing.T) {
	t.Parallel()
	// Arrange
	database := getDb(t)
	autoMigrate(t, database)

	service := &DBPlayerService{Database: database}

	creator := &domain.Creator{
		BaseObject: domain.BaseObject{ID: uuid.MustParse("84ac4166-7202-480b-93ff-5cab13514436")},
		Nickname:   "a",
		AuthID:     "a",
	}
	quiz := &domain.Quiz{Name: "test", Creator: creator}
	game := &domain.Game{Quiz: quiz}
	players := []*domain.Player{
		{Game: game, BaseObject: domain.BaseObject{ID: uuid.MustParse("3c97f06b-1078-46ef-a2c3-71fc4d9a3d3d")}, Nickname: "A"},
		{Game: game, BaseObject: domain.BaseObject{ID: uuid.MustParse("f065f531-1f62-4a19-9c6c-d8fe5ed518d4")}, Nickname: "B"},
	}

	database.CreateInBatches(players, 10)

	// Act
	result, err := service.GetByID(players[0].ID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "A", result.Nickname)
}

func TestDBPlayerService_GetByID_ReturnsDatabaseError(t *testing.T) {
	t.Parallel()
	// Arrange
	database := getDb(t)

	// By not running this, we're sure it will return an error
	// autoMigrate(t, database)

	service := &DBPlayerService{Database: database}

	// Act
	result, err := service.GetByID(uuid.MustParse("3c97f06b-1078-46ef-a2c3-71fc4d9a3d3d"))

	// Assert
	assert.Empty(t, result)
	assert.ErrorContains(t, err, "no such table")
}

func TestDBPlayerService_GetByGame_ReturnsAnyError(t *testing.T) {
	t.Parallel()
	// Arrange
	database := getDb(t)

	// By not running this, we're sure it will return an error
	// autoMigrate(t, database)

	service := &DBPlayerService{
		Database: database,
	}

	quizId := uuid.MustParse("238fe389-dede-4ee0-b26f-d2b1a65befac")

	// Act
	result, err := service.GetByGame(quizId)

	// Assert
	assert.Empty(t, result)
	assert.ErrorContains(t, err, "no such table")
}

func TestDBPlayerService_GetByGame_ReturnsExpectedGames(t *testing.T) {
	t.Parallel()
	// Arrange
	database := getDb(t)
	autoMigrate(t, database)

	service := &DBPlayerService{
		Database: database,
	}

	quizId := uuid.MustParse("238fe389-dede-4ee0-b26f-d2b1a65befac")
	creator := &domain.Creator{}
	quiz := &domain.Quiz{BaseObject: domain.BaseObject{ID: quizId}, Creator: creator}
	games := []*domain.Game{
		{Quiz: quiz, Code: "dje452", Players: []*domain.Player{{Nickname: "A"}, {Nickname: "B"}}},
		{Quiz: quiz, Code: "lei502", Players: []*domain.Player{{Nickname: "C"}, {Nickname: "D"}}},
	}

	if err := database.CreateInBatches(games, 10).Error; err != nil {
		t.Fatal(err)
	}

	// Act
	result, err := service.GetByGame(games[0].ID)

	// Assert
	assert.NoError(t, err)

	if assert.Len(t, result, 2) {
		assert.Equal(t, games[0].Players[0].Nickname, result[0].Nickname)
		assert.Equal(t, games[0].Players[1].Nickname, result[1].Nickname)
	}
}

func TestDBPlayerService_Create_CreatesItem(t *testing.T) {
	t.Parallel()
	// Arrange
	database := getDb(t)
	autoMigrate(t, database)

	service := &DBPlayerService{
		Database: database,
	}

	quizId := uuid.MustParse("238fe389-dede-4ee0-b26f-d2b1a65befac")
	creator := &domain.Creator{}
	quiz := &domain.Quiz{BaseObject: domain.BaseObject{ID: quizId}, Creator: creator}
	games := []*domain.Game{
		{Quiz: quiz, Code: "dje452"},
	}

	if err := database.CreateInBatches(games, 10).Error; err != nil {
		t.Fatal(err)
	}

	player := &domain.Player{GameID: games[0].ID}

	// Act
	err := service.Create(player)

	// Assert
	assert.NoError(t, err)

	var result *domain.Player
	if err := database.First(&result).Error; err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, player.ID, result.ID)
	assert.NotEmpty(t, result.Nickname)
}

func TestDBPlayerService_Create_ReturnsAnyError(t *testing.T) {
	t.Parallel()
	// Arrange
	database := getDb(t)

	// By not running this, we're sure it will return an error
	// autoMigrate(t, database)

	service := &DBPlayerService{
		Database: database,
	}

	// Act
	err := service.Create(&domain.Player{})

	// Assert
	assert.ErrorContains(t, err, "no such table")
}

func TestDBPlayerService_Delete_DeletesCorrectly(t *testing.T) {
	t.Parallel()
	// Arrange
	database := getDb(t)
	autoMigrate(t, database)

	service := &DBPlayerService{
		Database: database,
	}

	player := &domain.Player{
		BaseObject: domain.BaseObject{ID: uuid.MustParse("238fe389-dede-4ee0-b26f-d2b1a65befac")},
		Game: &domain.Game{
			Quiz: &domain.Quiz{
				Creator: &domain.Creator{},
			},
		},
	}
	database.Create(player)

	// Act
	err := service.Delete(player)

	// Assert
	assert.NoError(t, err)

	assert.ErrorContains(t, database.First(&domain.Player{}).Error, "not found")
}

func TestDBPlayerService_Delete_ReturnsAnyError(t *testing.T) {
	t.Parallel()
	// Arrange
	database := getDb(t)

	// By not running this, we're sure it will return an error
	// autoMigrate(t, database)

	service := &DBPlayerService{
		Database: database,
	}

	// Act
	err := service.Delete(&domain.Player{BaseObject: domain.BaseObject{ID: uuid.MustParse("238fe389-dede-4ee0-b26f-d2b1a65befac")}})

	// Assert
	assert.ErrorContains(t, err, "no such table")
}
