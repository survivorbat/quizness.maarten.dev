package domain

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGame_IsInProgress_ReturnsTrueOnStarted(t *testing.T) {
	t.Parallel()
	// Arrange
	game := &Game{StartTime: time.Now()}

	// Act
	result := game.IsInProgress()

	// Assert
	assert.True(t, result)
}

func TestGame_IsInProgress_ReturnsFalseOnNotStarted(t *testing.T) {
	t.Parallel()
	// Arrange
	game := &Game{}

	// Act
	result := game.IsInProgress()

	// Assert
	assert.False(t, result)
}
func TestGame_IsInProgress_ReturnsFalseOnStartedAndFinished(t *testing.T) {
	t.Parallel()
	// Arrange
	game := &Game{StartTime: time.Now(), FinishTime: time.Now()}

	// Act
	result := game.IsInProgress()

	// Assert
	assert.False(t, result)
}

func TestGame_Start_SetsStartTimeAndCode(t *testing.T) {
	t.Parallel()
	// Arrange
	game := Game{}

	// Act
	err := game.Start()

	// Assert
	assert.NoError(t, err)
	assert.False(t, game.StartTime.IsZero())
	assert.Len(t, game.Code, 6)
}

func TestGame_Start_ErrorsOnAlreadyStarted(t *testing.T) {
	t.Parallel()
	// Arrange
	game := Game{}

	_ = game.Start()

	// Act
	err := game.Start()

	// Assert
	assert.ErrorContains(t, err, "game has already started")
}

func TestGame_Finish_SetsFinishTime(t *testing.T) {
	t.Parallel()
	// Arrange
	game := Game{StartTime: time.Now()}

	// Act
	err := game.Finish()

	// Assert
	assert.NoError(t, err)
	assert.False(t, game.FinishTime.IsZero())
}

func TestGame_Finish_ErrorsOnNotStarted(t *testing.T) {
	t.Parallel()
	// Arrange
	game := Game{}

	// Act
	err := game.Finish()

	// Assert
	assert.ErrorContains(t, err, "game has not started")
}

func TestGame_Finish_ErrorsOnAlreadyFinished(t *testing.T) {
	t.Parallel()
	// Arrange
	game := Game{StartTime: time.Now()}

	_ = game.Finish()

	// Act
	err := game.Finish()

	// Assert
	assert.ErrorContains(t, err, "game has already finished")
}

func TestGame_AnswerQuestion_ReturnsErrorOnWrongQuestion(t *testing.T) {
	t.Parallel()
	// Arrange
	questionID := uuid.MustParse("ae2a9fd4-f861-4c99-a9bd-2bf49e1b1cd8")
	playerID := uuid.MustParse("a28ca8c6-63d9-45a2-b990-9b41e306f156")
	optionID := uuid.MustParse("c7ff1cdf-72d3-4ea9-ae48-e1c1d61f8bc8")

	game := &Game{CurrentQuestion: uuid.MustParse("f6d2fa67-c12d-4096-958b-18206fbf3538")}

	// Act
	answer, err := game.AnswerQuestion(playerID, questionID, optionID)

	// Assert
	assert.Nil(t, answer)
	assert.ErrorContains(t, err, "not the current question")
}

func TestGame_AnswerQuestion_ReturnsErrorOnDeadlinePassed(t *testing.T) {
	t.Parallel()
	// Arrange
	questionID := uuid.MustParse("ae2a9fd4-f861-4c99-a9bd-2bf49e1b1cd8")
	playerID := uuid.MustParse("a28ca8c6-63d9-45a2-b990-9b41e306f156")
	optionID := uuid.MustParse("c7ff1cdf-72d3-4ea9-ae48-e1c1d61f8bc8")

	game := &Game{CurrentQuestion: questionID, CurrentDeadline: time.Now()}

	// Act
	answer, err := game.AnswerQuestion(playerID, questionID, optionID)

	// Assert
	assert.Nil(t, answer)
	assert.ErrorContains(t, err, "deadline passed")
}

func TestGame_AnswerQuestion_ReturnsErrorOnPlayerNotInGame(t *testing.T) {
	t.Parallel()
	// Arrange
	questionID := uuid.MustParse("ae2a9fd4-f861-4c99-a9bd-2bf49e1b1cd8")
	playerID := uuid.MustParse("a28ca8c6-63d9-45a2-b990-9b41e306f156")
	optionID := uuid.MustParse("c7ff1cdf-72d3-4ea9-ae48-e1c1d61f8bc8")

	game := &Game{CurrentQuestion: questionID, CurrentDeadline: time.Now().Add(5 * time.Hour)}

	// Act
	answer, err := game.AnswerQuestion(playerID, questionID, optionID)

	// Assert
	assert.Nil(t, answer)
	assert.ErrorContains(t, err, "player not in game")
}

func TestGame_AnswerQuestion_ReturnsErrorOnAlreadyAnswered(t *testing.T) {
	t.Parallel()
	// Arrange
	questionID := uuid.MustParse("ae2a9fd4-f861-4c99-a9bd-2bf49e1b1cd8")
	playerID := uuid.MustParse("a28ca8c6-63d9-45a2-b990-9b41e306f156")
	optionID := uuid.MustParse("c7ff1cdf-72d3-4ea9-ae48-e1c1d61f8bc8")

	game := &Game{
		BaseObject:      BaseObject{ID: uuid.MustParse("bd787fed-8a2e-40d3-abf6-6b90fa89f862")},
		CurrentQuestion: questionID,
		CurrentDeadline: time.Now().Add(5 * time.Hour),
		Answers:         []*GameAnswer{{PlayerID: playerID, QuestionID: questionID}},
		Players:         []*Player{{BaseObject: BaseObject{ID: playerID}}},
	}

	// Act
	answer, err := game.AnswerQuestion(playerID, questionID, optionID)

	// Assert
	assert.Nil(t, answer)
	assert.ErrorContains(t, err, "player has already submitted an answer")
}

func TestGame_AnswerQuestion_ReturnsAnswer(t *testing.T) {
	t.Parallel()
	// Arrange
	questionID := uuid.MustParse("ae2a9fd4-f861-4c99-a9bd-2bf49e1b1cd8")
	playerID := uuid.MustParse("a28ca8c6-63d9-45a2-b990-9b41e306f156")
	optionID := uuid.MustParse("c7ff1cdf-72d3-4ea9-ae48-e1c1d61f8bc8")

	game := &Game{
		BaseObject:      BaseObject{ID: uuid.MustParse("bd787fed-8a2e-40d3-abf6-6b90fa89f862")},
		CurrentQuestion: questionID,
		CurrentDeadline: time.Now().Add(5 * time.Hour),
		Players:         []*Player{{BaseObject: BaseObject{ID: playerID}}},
	}

	// Act
	answer, err := game.AnswerQuestion(playerID, questionID, optionID)

	// Assert
	assert.NoError(t, err)

	assert.Equal(t, game.ID, answer.GameID)
	assert.Equal(t, playerID, answer.PlayerID)
	assert.Equal(t, optionID, answer.OptionID)
	assert.Equal(t, questionID, answer.QuestionID)
}

func TestGame_Next_ReturnsErrorOnNotInProgress(t *testing.T) {
	t.Parallel()
	// Arrange
	game := &Game{
		BaseObject: BaseObject{ID: uuid.MustParse("bd787fed-8a2e-40d3-abf6-6b90fa89f862")},
	}

	// Act
	err := game.Next()

	// Assert
	assert.ErrorContains(t, err, "game is not in progress")
}

func TestGame_Next_ReturnsErrorOnNotNoPlayers(t *testing.T) {
	t.Parallel()
	// Arrange
	game := &Game{
		BaseObject: BaseObject{ID: uuid.MustParse("bd787fed-8a2e-40d3-abf6-6b90fa89f862")},
		StartTime:  time.Now(),
	}

	// Act
	err := game.Next()

	// Assert
	assert.ErrorContains(t, err, "can only start with 2 or more players")
}

func TestGame_Next_ReturnsErrorOnNoMoreQuestions(t *testing.T) {
	t.Parallel()
	// Arrange
	game := &Game{
		BaseObject: BaseObject{ID: uuid.MustParse("bd787fed-8a2e-40d3-abf6-6b90fa89f862")},
		StartTime:  time.Now(),
		Quiz:       &Quiz{},
		Players:    Players{{}},
	}

	// Act
	err := game.Next()

	// Assert
	assert.ErrorContains(t, err, "can only start with 2 or more players")
}

func TestGame_Next_SetsNextQuestionOnNil(t *testing.T) {
	t.Parallel()
	// Arrange
	game := &Game{
		BaseObject: BaseObject{ID: uuid.MustParse("bd787fed-8a2e-40d3-abf6-6b90fa89f862")},
		StartTime:  time.Now(),
		Players:    Players{{}, {}},
		Quiz: &Quiz{
			MultipleChoiceQuestions: []*MultipleChoiceQuestion{
				{
					BaseQuestion: BaseQuestion{
						BaseObject: BaseObject{ID: uuid.MustParse("32acdba2-3472-4489-82a2-426c22ff529c")},
					},
				},
			},
		},
	}

	// Act
	err := game.Next()

	// Assert
	assert.NoError(t, err)

	assert.Equal(t, uuid.MustParse("32acdba2-3472-4489-82a2-426c22ff529c"), game.CurrentQuestion)
	assert.False(t, game.CurrentDeadline.IsZero())
}
