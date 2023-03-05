package services

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/survivorbat/qq.maarten.dev/server/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func TestQuizService_GetByCreator_ReturnsUser(t *testing.T) {
	t.Parallel()
	// Arrange
	database, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	_ = database.AutoMigrate(&domain.Quiz{}, &domain.Creator{})

	service := &QuizService{Database: database}

	creators := []*domain.Creator{
		{BaseObject: domain.BaseObject{ID: uuid.MustParse("3c97f06b-1078-46ef-a2c3-71fc4d9a3d3d")}, AuthID: "a"},
		{BaseObject: domain.BaseObject{ID: uuid.MustParse("5dff22a7-afc7-4e4b-a63d-9903dedd66bf")}, AuthID: "b"},
	}

	quizzes := []*domain.Quiz{
		{BaseObject: domain.BaseObject{ID: uuid.MustParse("6aacfb41-e478-46ec-857e-11221f2a97fc")}, Creator: creators[0]},
		{BaseObject: domain.BaseObject{ID: uuid.MustParse("6b362ab3-f164-4073-82a4-f7c3d2010947")}, Creator: creators[0]},
		{BaseObject: domain.BaseObject{ID: uuid.MustParse("e7c3d165-0419-4e4d-beb7-45ea4e0e908e")}, Creator: creators[1]},
	}

	database.CreateInBatches(quizzes, 10)

	// Act
	result, err := service.GetByCreator(quizzes[0].CreatorID)

	// Assert
	assert.NoError(t, err)

	if assert.Len(t, result, 2) {
		assert.Equal(t, quizzes[0].ID, result[0].ID)
		assert.Equal(t, quizzes[1].ID, result[1].ID)
	}
}

func TestQuizService_GetByCreator_ReturnsDatabaseError(t *testing.T) {
	t.Parallel()
	// Arrange
	database, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	// By not running this, we're sure it will return an error
	//_ = database.AutoMigrate(&domain.Quiz{})

	service := &QuizService{Database: database}

	// Act
	result, err := service.GetByCreator(uuid.MustParse("3c97f06b-1078-46ef-a2c3-71fc4d9a3d3d"))

	// Assert
	assert.Empty(t, result)
	assert.ErrorContains(t, err, "no such table")
}
