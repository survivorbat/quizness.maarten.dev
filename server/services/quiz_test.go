package services

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/survivorbat/qq.maarten.dev/server/domain"
	"gorm.io/gorm"
	"testing"
)

func TestDBQuizService_GetByID_ReturnsExpected(t *testing.T) {
	t.Parallel()
	// Arrange
	database := getDb(t)
	autoMigrate(t, database)

	service := &DBQuizService{Database: database}

	creators := []*domain.Creator{
		{BaseObject: domain.BaseObject{ID: uuid.MustParse("3c97f06b-1078-46ef-a2c3-71fc4d9a3d3d")}, Nickname: "a", AuthID: "a"},
		{BaseObject: domain.BaseObject{ID: uuid.MustParse("5dff22a7-afc7-4e4b-a63d-9903dedd66bf")}, Nickname: "b", AuthID: "b"},
	}

	quizzes := []*domain.Quiz{
		{BaseObject: domain.BaseObject{ID: uuid.MustParse("6aacfb41-e478-46ec-857e-11221f2a97fc")}, Creator: creators[0], Games: []*domain.Game{{Code: "abc"}}},
		{BaseObject: domain.BaseObject{ID: uuid.MustParse("6b362ab3-f164-4073-82a4-f7c3d2010947")}, Creator: creators[0], Games: []*domain.Game{{Code: "def"}}},
		{BaseObject: domain.BaseObject{ID: uuid.MustParse("e7c3d165-0419-4e4d-beb7-45ea4e0e908e")}, Creator: creators[1], Games: []*domain.Game{{Code: "ghi"}}},
	}

	database.CreateInBatches(quizzes, 10)

	// Act
	result, err := service.GetByID(quizzes[0].ID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, quizzes[0].ID, result.ID)
	assert.Equal(t, quizzes[0].Games[0].Code, result.Games[0].Code)
}

func TestDBQuizService_GetByID_ReturnsDatabaseError(t *testing.T) {
	t.Parallel()
	// Arrange
	database := getDb(t)

	// By not running this, we're sure it will return an error
	// autoMigrate(t, database)

	service := &DBQuizService{Database: database}

	// Act
	result, err := service.GetByID(uuid.MustParse("3c97f06b-1078-46ef-a2c3-71fc4d9a3d3d"))

	// Assert
	assert.Empty(t, result)
	assert.ErrorContains(t, err, "no such table")
}

func TestDBQuizService_GetByCreator_ReturnsUser(t *testing.T) {
	t.Parallel()
	// Arrange
	database := getDb(t)
	autoMigrate(t, database)

	service := &DBQuizService{Database: database}

	creators := []*domain.Creator{
		{BaseObject: domain.BaseObject{ID: uuid.MustParse("3c97f06b-1078-46ef-a2c3-71fc4d9a3d3d")}, Nickname: "a", AuthID: "a"},
		{BaseObject: domain.BaseObject{ID: uuid.MustParse("5dff22a7-afc7-4e4b-a63d-9903dedd66bf")}, Nickname: "b", AuthID: "b"},
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

func TestDBQuizService_GetByCreator_ReturnsDatabaseError(t *testing.T) {
	t.Parallel()
	// Arrange
	database := getDb(t)

	// By not running this, we're sure it will return an error
	// autoMigrate(t, database)

	service := &DBQuizService{Database: database}

	// Act
	result, err := service.GetByCreator(uuid.MustParse("3c97f06b-1078-46ef-a2c3-71fc4d9a3d3d"))

	// Assert
	assert.Empty(t, result)
	assert.ErrorContains(t, err, "no such table")
}

func TestDBQuizService_CreateOrUpdate_CreatesNewQuiz(t *testing.T) {
	t.Parallel()
	// Arrange
	database := getDb(t)
	autoMigrate(t, database)

	service := &DBQuizService{Database: database}

	quiz := &domain.Quiz{
		Name:    "test",
		Creator: &domain.Creator{},
		MultipleChoiceQuestions: []*domain.MultipleChoiceQuestion{
			{
				BaseQuestion: domain.BaseQuestion{
					Title: "abc",
				},
				Options: []*domain.QuestionOption{
					{TextOption: "def"},
				},
			},
		},
	}

	// Act
	err := service.CreateOrUpdate(quiz)

	// Assert
	assert.NoError(t, err)

	var result *domain.Quiz
	if err := database.First(&result).Error; err != nil {
		t.Fatal(err.Error())
	}
	assert.Equal(t, "test", result.Name)
}

func TestDBQuizService_CreateOrUpdate_UpdatesExisting(t *testing.T) {
	t.Parallel()
	// Arrange
	database := getDb(t)
	autoMigrate(t, database)

	service := &DBQuizService{Database: database}

	existing := &domain.Quiz{
		Name:    "old",
		Creator: &domain.Creator{Nickname: "abc", AuthID: "def"},
		MultipleChoiceQuestions: []*domain.MultipleChoiceQuestion{
			{
				BaseQuestion: domain.BaseQuestion{
					Title: "abc",
				},
				Options: []*domain.QuestionOption{
					{TextOption: "def"},
				},
			},
		},
	}
	if err := database.Create(existing).Error; err != nil {
		t.Fatal(err)
	}

	update := &domain.Quiz{
		BaseObject: domain.BaseObject{ID: existing.ID},
		Name:       "new",
		Creator:    &domain.Creator{Nickname: "abc", AuthID: "def"},
		MultipleChoiceQuestions: []*domain.MultipleChoiceQuestion{
			{
				BaseQuestion: domain.BaseQuestion{
					Title: "def",
				},
				Options: []*domain.QuestionOption{
					{TextOption: "abc"},
				},
			},
		},
	}

	// Act
	err := service.CreateOrUpdate(update)

	// Assert
	assert.NoError(t, err)

	var result *domain.Quiz
	if err := database.Preload("MultipleChoiceQuestions.Options").First(&result).Error; err != nil {
		t.Fatal(err.Error())
	}

	assert.Equal(t, "new", result.Name)
	assert.Equal(t, "def", result.MultipleChoiceQuestions[0].Title)
	assert.Equal(t, "abc", result.MultipleChoiceQuestions[0].Options[0].TextOption)
}

func TestDBQuizService_CreateOrUpdate_ReturnsDatabaseError(t *testing.T) {
	t.Parallel()
	// Arrange
	database := getDb(t)

	// By not running this, we're sure it will return an error
	// autoMigrate(t, database)

	service := &DBQuizService{Database: database}

	// Act
	err := service.CreateOrUpdate(&domain.Quiz{})

	// Assert
	assert.ErrorContains(t, err, "no such table")
}

func TestDBQuizService_Delete_DeletesQuiz(t *testing.T) {
	t.Parallel()
	// Arrange
	database := getDb(t)
	autoMigrate(t, database)

	service := &DBQuizService{Database: database}

	quiz := &domain.Quiz{
		Name:    "test",
		Creator: &domain.Creator{},
		MultipleChoiceQuestions: []*domain.MultipleChoiceQuestion{
			{
				BaseQuestion: domain.BaseQuestion{
					Title: "abc",
				},
				Options: []*domain.QuestionOption{
					{TextOption: "def"},
				},
			},
		},
	}

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

	var questions []*domain.MultipleChoiceQuestion
	if err := database.Find(&questions).Error; err != nil {
		t.Fatal(err.Error())
	}
	assert.Len(t, questions, 0)

	var options []*domain.QuestionOption
	if err := database.Find(&options).Error; err != nil {
		t.Fatal(err.Error())
	}
	assert.Len(t, options, 0)
}

func TestDBQuizService_Delete_ReturnsDatabaseError(t *testing.T) {
	t.Parallel()
	// Arrange
	database := getDb(t)

	// By not running this, we're sure it will return an error
	// autoMigrate(t, database)

	service := &DBQuizService{Database: database}

	// Act
	err := service.Delete(uuid.MustParse("3c97f06b-1078-46ef-a2c3-71fc4d9a3d3d"))

	// Assert
	assert.ErrorContains(t, err, "no such table")
}
