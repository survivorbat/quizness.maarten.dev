package services

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/survivorbat/qq.maarten.dev/server/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func TestQuizService_GetByID_ReturnsUser(t *testing.T) {
	t.Parallel()
	// Arrange
	database, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	_ = database.AutoMigrate(&domain.Quiz{}, &domain.Creator{}, &domain.MultipleChoiceQuestion{}, &domain.QuestionOption{})

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
	result, err := service.GetByID(quizzes[0].ID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, quizzes[0].ID, result.ID)
}

func TestQuizService_GetByID_ReturnsDatabaseError(t *testing.T) {
	t.Parallel()
	// Arrange
	database, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	// By not running this, we're sure it will return an error
	//_ = database.AutoMigrate(&domain.Quiz{})

	service := &QuizService{Database: database}

	// Act
	result, err := service.GetByID(uuid.MustParse("3c97f06b-1078-46ef-a2c3-71fc4d9a3d3d"))

	// Assert
	assert.Empty(t, result)
	assert.ErrorContains(t, err, "no such table")
}

func TestQuizService_GetByCreator_ReturnsUser(t *testing.T) {
	t.Parallel()
	// Arrange
	database, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	_ = database.AutoMigrate(&domain.Quiz{}, &domain.Creator{}, &domain.MultipleChoiceQuestion{}, &domain.QuestionOption{})

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

func TestQuizService_Create_CreatesQuiz(t *testing.T) {
	t.Parallel()
	// Arrange
	database, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	_ = database.AutoMigrate(&domain.Quiz{})

	service := &QuizService{Database: database}

	// Act
	err := service.Create(&domain.Quiz{Name: "test"})

	// Assert
	assert.NoError(t, err)

	var result *domain.Quiz
	if err := database.First(&result).Error; err != nil {
		t.Fatal(err.Error())
	}
	assert.Equal(t, "test", result.Name)
}

func TestQuizService_Create_ReturnsDatabaseError(t *testing.T) {
	t.Parallel()
	// Arrange
	database, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	// By not running this, we're sure it will return an error
	//_ = database.AutoMigrate(&domain.Quiz{})

	service := &QuizService{Database: database}

	// Act
	err := service.Create(&domain.Quiz{})

	// Assert
	assert.ErrorContains(t, err, "no such table")
}

func TestQuizService_Delete_DeletesQuiz(t *testing.T) {
	t.Parallel()
	// Arrange
	database, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	_ = database.AutoMigrate(&domain.Quiz{})

	service := &QuizService{Database: database}

	quiz := &domain.Quiz{Name: "test"}
	if err := database.Create(quiz).Error; err != nil {
		t.Fatal(err.Error())
	}

	// Act
	err := service.Delete(quiz.ID)

	// Assert
	assert.NoError(t, err)

	var result *domain.Quiz
	err = database.First(&result).Error
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestQuizService_Delete_ReturnsDatabaseError(t *testing.T) {
	t.Parallel()
	// Arrange
	database, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	// By not running this, we're sure it will return an error
	//_ = database.AutoMigrate(&domain.Quiz{})

	service := &QuizService{Database: database}

	// Act
	err := service.Delete(uuid.MustParse("3c97f06b-1078-46ef-a2c3-71fc4d9a3d3d"))

	// Assert
	assert.ErrorContains(t, err, "no such table")
}
