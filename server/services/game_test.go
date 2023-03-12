package services

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/survivorbat/qq.maarten.dev/server/domain"
	"testing"
)

func TestDBGameService_GetByQuiz_ReturnsAnyError(t *testing.T) {
	t.Parallel()
	// Arrange
	database := getDb(t)

	// By not running this, we're sure it will return an error
	// autoMigrate(t, database)

	service := &DBGameService{
		Database: database,
	}

	quizId := uuid.MustParse("238fe389-dede-4ee0-b26f-d2b1a65befac")

	// Act
	result, err := service.GetByQuiz(quizId)

	// Assert
	assert.Empty(t, result)
	assert.ErrorContains(t, err, "no such table")
}

func TestDBGameService_GetByQuiz_ReturnsExpectedGames(t *testing.T) {
	t.Parallel()
	// Arrange
	database := getDb(t)
	autoMigrate(t, database)

	service := &DBGameService{
		Database: database,
	}

	quizId := uuid.MustParse("238fe389-dede-4ee0-b26f-d2b1a65befac")
	creator := &domain.Creator{}
	quizA := &domain.Quiz{BaseObject: domain.BaseObject{ID: quizId}, Creator: creator}
	quizB := &domain.Quiz{Creator: creator}
	games := []*domain.Game{
		{Quiz: quizA, Code: "dje452"},
		{Quiz: quizA, Code: "lei502"},
		{Quiz: quizB, Code: "gde235"},
	}

	if err := database.CreateInBatches(games, 10).Error; err != nil {
		t.Fatal(err)
	}

	// Act
	result, err := service.GetByQuiz(quizId)

	// Assert
	assert.NoError(t, err)

	if assert.Len(t, result, 2) {
		assert.Equal(t, games[0].Code, result[0].Code)
		assert.Equal(t, games[1].Code, result[1].Code)
	}
}

func TestDBGameService_Create_ReturnsAnyError(t *testing.T) {
	t.Parallel()
	// Arrange
	database := getDb(t)

	// By not running this, we're sure it will return an error
	// autoMigrate(t, database)

	service := &DBGameService{
		Database: database,
	}

	// Act
	err := service.Create(&domain.Game{})

	// Assert
	assert.ErrorContains(t, err, "no such table")
}

func TestDBGameService_Create_CreatesGame(t *testing.T) {
	t.Parallel()
	// Arrange
	database := getDb(t)
	autoMigrate(t, database)

	service := &DBGameService{
		Database: database,
	}

	quizId := uuid.MustParse("238fe389-dede-4ee0-b26f-d2b1a65befac")
	creator := &domain.Creator{}
	quiz := &domain.Quiz{BaseObject: domain.BaseObject{ID: quizId}, Creator: creator}

	if err := database.Create(quiz).Error; err != nil {
		t.Fatal(err)
	}

	// Act
	err := service.Create(&domain.Game{Quiz: quiz, PlayerLimit: 20})

	// Assert
	assert.NoError(t, err)

	var result *domain.Game
	if err := database.First(&result).Error; err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, uint(20), result.PlayerLimit)
}
