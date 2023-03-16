package services

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/survivorbat/qq.maarten.dev/server/domain"
	"testing"
	"time"
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

func TestDBGameService_GetByID_ReturnsExpected(t *testing.T) {
	t.Parallel()
	// Arrange
	database := getDb(t)
	autoMigrate(t, database)

	service := &DBGameService{Database: database}

	creator := &domain.Creator{
		BaseObject: domain.BaseObject{ID: uuid.MustParse("3c97f06b-1078-46ef-a2c3-71fc4d9a3d3d")},
		Nickname:   "a",
		AuthID:     "a",
	}
	quiz := &domain.Quiz{Name: "test", Creator: creator}

	games := []*domain.Game{
		{BaseObject: domain.BaseObject{ID: uuid.MustParse("6aacfb41-e478-46ec-857e-11221f2a97fc")}, Quiz: quiz, Code: "A2DFGH"},
		{BaseObject: domain.BaseObject{ID: uuid.MustParse("9180d979-2de5-4df2-a6ee-07eec2f79d92")}, Quiz: quiz, Code: "920LEK"},
	}

	database.CreateInBatches(games, 10)

	// Act
	result, err := service.GetByID(games[0].ID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, games[0].ID, result.ID)
	assert.Equal(t, games[0].Code, result.Code)
	assert.Equal(t, quiz.Name, result.Quiz.Name)
}

func TestDBGameService_GetByID_ReturnsDatabaseError(t *testing.T) {
	t.Parallel()
	// Arrange
	database := getDb(t)

	// By not running this, we're sure it will return an error
	// autoMigrate(t, database)

	service := &DBGameService{Database: database}

	// Act
	result, err := service.GetByID(uuid.MustParse("3c97f06b-1078-46ef-a2c3-71fc4d9a3d3d"))

	// Assert
	assert.Empty(t, result)
	assert.ErrorContains(t, err, "no such table")
}

func TestDBGameService_GetByCode_ReturnsExpected(t *testing.T) {
	t.Parallel()
	// Arrange
	database := getDb(t)
	autoMigrate(t, database)

	service := &DBGameService{Database: database}

	creator := &domain.Creator{
		BaseObject: domain.BaseObject{ID: uuid.MustParse("3c97f06b-1078-46ef-a2c3-71fc4d9a3d3d")},
		Nickname:   "a",
		AuthID:     "a",
	}
	quiz := &domain.Quiz{Name: "test", Creator: creator}

	games := []*domain.Game{
		{BaseObject: domain.BaseObject{ID: uuid.MustParse("6aacfb41-e478-46ec-857e-11221f2a97fc")}, Quiz: quiz, Code: "A2DFGH"},
		{BaseObject: domain.BaseObject{ID: uuid.MustParse("9180d979-2de5-4df2-a6ee-07eec2f79d92")}, Quiz: quiz, Code: "920LEK"},
	}

	database.CreateInBatches(games, 10)

	// Act
	result, err := service.GetByCode(games[0].Code)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, games[0].ID, result.ID)
	assert.Equal(t, games[0].Code, result.Code)
}

func TestDBGameService_GetByCode_ReturnsDatabaseError(t *testing.T) {
	t.Parallel()
	// Arrange
	database := getDb(t)

	// By not running this, we're sure it will return an error
	// autoMigrate(t, database)

	service := &DBGameService{Database: database}

	// Act
	result, err := service.GetByCode("A2DFGH")

	// Assert
	assert.Empty(t, result)
	assert.ErrorContains(t, err, "no such table")
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

func TestDBGameService_Start_ReturnsErrorIfStarted(t *testing.T) {
	t.Parallel()
	// Arrange
	database := getDb(t)
	autoMigrate(t, database)

	service := &DBGameService{
		Database: database,
	}

	game := &domain.Game{
		BaseObject: domain.BaseObject{ID: uuid.MustParse("238fe389-dede-4ee0-b26f-d2b1a65befac")},
		StartTime:  time.Now(),
		Quiz: &domain.Quiz{
			Creator: &domain.Creator{},
		},
	}
	database.Create(game)

	// Act
	err := service.Start(game)

	// Assert
	assert.ErrorContains(t, err, "game has already started")
}

func TestDBGameService_Start_StartsGame(t *testing.T) {
	t.Parallel()
	// Arrange
	database := getDb(t)
	autoMigrate(t, database)

	service := &DBGameService{
		Database: database,
	}

	game := &domain.Game{
		BaseObject: domain.BaseObject{ID: uuid.MustParse("238fe389-dede-4ee0-b26f-d2b1a65befac")},
		Quiz: &domain.Quiz{
			Creator:                 &domain.Creator{},
			MultipleChoiceQuestions: []*domain.MultipleChoiceQuestion{{}, {}},
		},
	}
	database.Create(game)

	// Act
	err := service.Start(game)

	// Assert
	assert.NoError(t, err)

	var result *domain.Game
	if err := database.First(&result).Error; err != nil {
		t.Fatal(err)
	}
	assert.False(t, result.StartTime.IsZero())
}

func TestDBGameService_Finish_ReturnsErrorIfAlreadyFinished(t *testing.T) {
	t.Parallel()
	// Arrange
	database := getDb(t)
	autoMigrate(t, database)

	service := &DBGameService{
		Database: database,
	}

	game := &domain.Game{
		BaseObject: domain.BaseObject{ID: uuid.MustParse("238fe389-dede-4ee0-b26f-d2b1a65befac")},
		StartTime:  time.Now(),
		FinishTime: time.Now(),
		Quiz: &domain.Quiz{
			Creator: &domain.Creator{},
		},
	}
	database.Create(game)

	// Act
	err := service.Finish(game)

	// Assert
	assert.ErrorContains(t, err, "game has already finished")
}

func TestDBGameService_Finish_FinishesGame(t *testing.T) {
	t.Parallel()
	// Arrange
	database := getDb(t)
	autoMigrate(t, database)

	service := &DBGameService{
		Database: database,
	}

	game := &domain.Game{
		BaseObject: domain.BaseObject{ID: uuid.MustParse("238fe389-dede-4ee0-b26f-d2b1a65befac")},
		StartTime:  time.Now(),
		Quiz: &domain.Quiz{
			Creator: &domain.Creator{},
		},
	}
	database.Create(game)

	// Act
	err := service.Finish(game)

	// Assert
	assert.NoError(t, err)

	var result *domain.Game
	if err := database.First(&result).Error; err != nil {
		t.Fatal(err)
	}
	assert.False(t, result.FinishTime.IsZero())
}

func TestDBGameService_Next_StartsNextQuestion(t *testing.T) {
	t.Parallel()
	// Arrange
	database := getDb(t)
	autoMigrate(t, database)

	service := &DBGameService{
		Database: database,
	}

	game := &domain.Game{
		BaseObject: domain.BaseObject{ID: uuid.MustParse("238fe389-dede-4ee0-b26f-d2b1a65befac")},
		StartTime:  time.Now(),
		Players:    []*domain.Player{{}, {}},
		Quiz: &domain.Quiz{
			Creator: &domain.Creator{},
			MultipleChoiceQuestions: []*domain.MultipleChoiceQuestion{
				{
					BaseQuestion: domain.BaseQuestion{
						BaseObject: domain.BaseObject{ID: uuid.MustParse("c275bf4e-c839-495d-af9c-4f95d8dc05a5")},
						Order:      0,
					},
				},
				{
					BaseQuestion: domain.BaseQuestion{
						BaseObject: domain.BaseObject{ID: uuid.MustParse("ce454f0d-9d9c-4e39-bd86-3484a7283eec")},
						Order:      1,
					},
				},
			},
		},
	}
	database.Create(game)

	// Act
	err := service.Next(game)

	// Assert
	assert.NoError(t, err)

	var result *domain.Game
	if err := database.First(&result).Error; err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, uuid.MustParse("c275bf4e-c839-495d-af9c-4f95d8dc05a5"), result.CurrentQuestion)
}

func TestDBGameService_Next_ReturnsErrorIfNotStarted(t *testing.T) {
	t.Parallel()
	// Arrange
	database := getDb(t)
	autoMigrate(t, database)

	service := &DBGameService{
		Database: database,
	}

	game := &domain.Game{
		BaseObject: domain.BaseObject{ID: uuid.MustParse("238fe389-dede-4ee0-b26f-d2b1a65befac")},
		Quiz: &domain.Quiz{
			Creator: &domain.Creator{},
		},
	}
	database.Create(game)

	// Act
	err := service.Next(game)

	// Assert
	assert.ErrorContains(t, err, "game is not in progress")
}
func TestDBGameService_AnswerQuestion_StartsAnswerQuestionQuestion(t *testing.T) {
	t.Parallel()
	// Arrange
	database := getDb(t)
	autoMigrate(t, database)

	questionId := uuid.MustParse("c275bf4e-c839-495d-af9c-4f95d8dc05a5")
	playerId := uuid.MustParse("62750588-5575-4a31-9cdf-2ffed23c7a15")
	optionId := uuid.MustParse("ecbffee9-c66a-4d33-9cdc-ac0e15da2982")

	service := &DBGameService{
		Database: database,
	}

	game := &domain.Game{
		BaseObject:      domain.BaseObject{ID: uuid.MustParse("238fe389-dede-4ee0-b26f-d2b1a65befac")},
		StartTime:       time.Now(),
		Players:         []*domain.Player{{BaseObject: domain.BaseObject{ID: playerId}}, {}},
		CurrentQuestion: questionId,
		CurrentDeadline: time.Now().Add(20 * time.Hour),
		Quiz: &domain.Quiz{
			Creator: &domain.Creator{},
			MultipleChoiceQuestions: []*domain.MultipleChoiceQuestion{
				{
					BaseQuestion: domain.BaseQuestion{
						BaseObject: domain.BaseObject{ID: questionId},
						Order:      0,
					},
				},
				{
					BaseQuestion: domain.BaseQuestion{
						BaseObject: domain.BaseObject{ID: uuid.MustParse("ce454f0d-9d9c-4e39-bd86-3484a7283eec")},
						Order:      1,
					},
				},
			},
		},
	}
	database.Create(game)

	// Act
	err := service.AnswerQuestion(game, questionId, playerId, optionId)

	// Assert
	assert.NoError(t, err)

	var result *domain.GameAnswer
	if err := database.First(&result).Error; err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, optionId, result.OptionID)
	assert.Equal(t, questionId, result.QuestionID)
	assert.Equal(t, playerId, result.PlayerID)
	assert.Equal(t, game.ID, result.GameID)
}

func TestDBGameService_AnswerQuestion_ReturnsErrorIfNotTheCurrentQuestion(t *testing.T) {
	t.Parallel()
	// Arrange
	database := getDb(t)
	autoMigrate(t, database)

	questionId := uuid.MustParse("c275bf4e-c839-495d-af9c-4f95d8dc05a5")
	playerId := uuid.MustParse("62750588-5575-4a31-9cdf-2ffed23c7a15")
	optionId := uuid.MustParse("ecbffee9-c66a-4d33-9cdc-ac0e15da2982")

	service := &DBGameService{
		Database: database,
	}

	game := &domain.Game{
		BaseObject: domain.BaseObject{ID: uuid.MustParse("238fe389-dede-4ee0-b26f-d2b1a65befac")},
		Quiz: &domain.Quiz{
			Creator: &domain.Creator{},
		},
	}
	database.Create(game)

	// Act
	err := service.AnswerQuestion(game, questionId, playerId, optionId)

	// Assert
	assert.ErrorContains(t, err, "not the current question")
}

func TestDBGameService_Delete_ReturnsErrorOnInProgress(t *testing.T) {
	t.Parallel()
	// Arrange
	database := getDb(t)
	autoMigrate(t, database)

	service := &DBGameService{
		Database: database,
	}

	game := &domain.Game{
		BaseObject: domain.BaseObject{ID: uuid.MustParse("238fe389-dede-4ee0-b26f-d2b1a65befac")},
		StartTime:  time.Now(),
		Quiz: &domain.Quiz{
			Creator: &domain.Creator{},
		},
	}
	database.Create(game)

	// Act
	err := service.Delete(game)

	// Assert
	assert.ErrorContains(t, err, "in progress")
}

func TestDBGameService_Delete_DeletesCorrectly(t *testing.T) {
	t.Parallel()
	// Arrange
	database := getDb(t)
	autoMigrate(t, database)

	service := &DBGameService{
		Database: database,
	}

	game := &domain.Game{
		BaseObject: domain.BaseObject{ID: uuid.MustParse("238fe389-dede-4ee0-b26f-d2b1a65befac")},
		Quiz: &domain.Quiz{
			Creator: &domain.Creator{},
		},
	}
	database.Create(game)

	// Act
	err := service.Delete(game)

	// Assert
	assert.NoError(t, err)

	var result *domain.Game
	assert.ErrorContains(t, database.First(&result).Error, "not found")
}
