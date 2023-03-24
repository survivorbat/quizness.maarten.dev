package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/ing-bank/gintestutil"
	"github.com/ing-bank/gormtestutil"
	"github.com/stretchr/testify/assert"
	"github.com/survivorbat/qq.maarten.dev/server/coordinator"
	"github.com/survivorbat/qq.maarten.dev/server/domain"
	"github.com/survivorbat/qq.maarten.dev/server/inputs"
	"github.com/survivorbat/qq.maarten.dev/server/routes/outputs"
	"golang.org/x/oauth2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

// Functions

func populateDatabase[T any](t *testing.T, database *gorm.DB, data ...T) {
	t.Helper()
	if err := database.Model(new(T)).CreateInBatches(data, 100).Error; err != nil {
		t.Fatal(err.Error())
	}
}

func performRequest(method string, server string, path string, auth string, body any) (*http.Response, error) {
	requestUrl, _ := url.Parse(fmt.Sprintf("%s/%s", server, path))
	request := &http.Request{
		Method: method,
		URL:    requestUrl,
	}

	if auth != "" {
		request.Header = map[string][]string{"Authorization": {fmt.Sprintf("Bearer %s", auth)}}
	}

	if body != nil {
		body, _ := json.Marshal(body)
		request.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	return http.DefaultClient.Do(request)
}

func getValue[T any, K any](t *testing.T, res *http.Response, err error, getKey func(T) K) K {
	t.Helper()

	if err != nil {
		t.Fatal(err)
	}

	var result T
	gintestutil.Response(t, &result, http.StatusOK, res)
	return getKey(result)
}

func getCreator(id uuid.UUID) *domain.Creator {
	result := &domain.Creator{BaseObject: domain.BaseObject{ID: id}, AuthID: id.String()}
	result.GenerateNickname()
	result.GenerateColors()
	return result
}

// Tests

func TestNewServer_GetQuizzes_ReturnsData(t *testing.T) {
	// Arrange
	instance := &Server{jwtSecret: "abc", oAuthConfig: &oauth2.Config{ClientID: "abc", ClientSecret: "abc", RedirectURL: "abc"}}
	instance.database = gormtestutil.NewMemoryDatabase(t, gormtestutil.WithName(t.Name()))

	// Test http server
	engine := gin.Default()
	_ = instance.Configure(engine)
	ts := httptest.NewServer(engine)

	userID := uuid.MustParse("7d87bab0-cf2d-45ae-bced-1de22db21a77")
	token, _ := instance.jwtService.GenerateToken(userID.String())

	quizzes := []*domain.Quiz{
		{Name: "def", Creator: getCreator(uuid.MustParse("dc0057c9-553d-40aa-a0bf-6fb98990c634"))},
		{
			Name:    "abc",
			Creator: getCreator(userID),
			MultipleChoiceQuestions: []*domain.MultipleChoiceQuestion{
				{
					BaseQuestion: domain.BaseQuestion{Title: "a"},
					Options: []*domain.QuestionOption{
						{TextOption: "abc"},
					},
				},
			},
		},
	}

	// Populate database
	populateDatabase(t, instance.database.Debug(), quizzes...)

	// Close it in the end
	defer ts.Close()

	// Act
	response, err := performRequest(http.MethodGet, ts.URL, "api/v1/quizzes", token, nil)

	// Assert
	assert.NoError(t, err)
	if !assert.NotNil(t, response) {
		t.FailNow()
	}

	assert.Equal(t, http.StatusOK, response.StatusCode)

	var result []*domain.Quiz
	ok := gintestutil.Response(t, &result, http.StatusOK, response)

	if ok && assert.Len(t, result, 1) {
		assert.Equal(t, quizzes[1].CreatorID, result[0].CreatorID)
		assert.Equal(t, quizzes[1].Name, result[0].Name)
		assert.Equal(t, quizzes[1].MultipleChoiceQuestions[0].Title, result[0].MultipleChoiceQuestions[0].Title)
		assert.Equal(t, quizzes[1].MultipleChoiceQuestions[0].Options[0].TextOption, result[0].MultipleChoiceQuestions[0].Options[0].TextOption)
	}
}

func TestNewServer_PostQuiz_ReturnsValidationErrors(t *testing.T) {
	tests := map[string]struct {
		input *inputs.Quiz
	}{
		"empty quiz": {
			input: &inputs.Quiz{},
		},
		"no questions": {
			input: &inputs.Quiz{
				Name:        "abc",
				Description: "def",
			},
		},
		"no options": {
			input: &inputs.Quiz{
				Name:        "abc",
				Description: "bcd",
				MultipleChoiceQuestions: []*inputs.MultipleChoiceQuestion{
					{
						Title:             "cde",
						Description:       "def",
						DurationInSeconds: 15,
						Category:          "egh",
						Order:             0,
					},
				},
			},
		},
		"wrong order": {
			input: &inputs.Quiz{
				Name:        "abc",
				Description: "bcd",
				MultipleChoiceQuestions: []*inputs.MultipleChoiceQuestion{
					{
						Title:             "cde",
						Description:       "def",
						DurationInSeconds: 15,
						Category:          "egh",
						Order:             0,
						Options: []*inputs.QuestionOption{
							{TextOption: "fgh"}, {TextOption: "ghi"},
							{TextOption: "hij"}, {TextOption: "ijk", Answer: true},
						},
					},
					{
						Title:             "cde",
						Description:       "def",
						DurationInSeconds: 15,
						Category:          "egh",
						Order:             0,
						Options: []*inputs.QuestionOption{
							{TextOption: "fgh", Answer: true}, {TextOption: "ghi"},
							{TextOption: "hij"}, {TextOption: "ijk"},
						},
					},
				},
			},
		},
		"missing answer": {
			input: &inputs.Quiz{
				Name:        "abc",
				Description: "bcd",
				MultipleChoiceQuestions: []*inputs.MultipleChoiceQuestion{
					{
						Title:             "cde",
						Description:       "def",
						DurationInSeconds: 15, Category: "egh",
						Order: 0,
						Options: []*inputs.QuestionOption{
							{TextOption: "hij"}, {TextOption: "ijk"},
						},
					},
				},
			},
		},
	}

	for name, testData := range tests {
		t.Run(name, func(t *testing.T) {
			// Arrange
			databaseOpen = sqlite.Open
			connection := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
			instance, _ := NewServer(connection, "abc", "abc", "abc", "abc")
			instance.database = gormtestutil.NewMemoryDatabase(t, gormtestutil.WithName(t.Name()))

			// Test http server
			engine := gin.Default()
			_ = instance.Configure(engine)
			ts := httptest.NewServer(engine)

			userID := uuid.MustParse("7d87bab0-cf2d-45ae-bced-1de22db21a77")
			token, _ := instance.jwtService.GenerateToken(userID.String())

			// Close it in the end
			defer ts.Close()

			// Act
			response, err := performRequest(http.MethodPost, ts.URL, "api/v1/quizzes", token, testData.input)

			// Assert
			assert.NoError(t, err)
			if !assert.NotNil(t, response) {
				t.FailNow()
			}

			assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		})
	}
}

func TestNewServer_PostQuiz_SavesQuiz(t *testing.T) {
	// Arrange
	instance := &Server{jwtSecret: "abc", oAuthConfig: &oauth2.Config{ClientID: "abc", ClientSecret: "abc", RedirectURL: "abc"}}
	instance.database = gormtestutil.NewMemoryDatabase(t, gormtestutil.WithName(t.Name()))

	// Test http server
	engine := gin.Default()
	_ = instance.Configure(engine)
	ts := httptest.NewServer(engine)

	userID := uuid.MustParse("7d87bab0-cf2d-45ae-bced-1de22db21a77")
	token, _ := instance.jwtService.GenerateToken(userID.String())

	// Close it in the end
	defer ts.Close()

	input := &inputs.Quiz{
		Name:        "abc",
		Description: "bcd",
		MultipleChoiceQuestions: []*inputs.MultipleChoiceQuestion{
			{
				Title:             "cde",
				Description:       "def",
				DurationInSeconds: 15,
				Category:          "egh",
				Order:             0,
				Options: []*inputs.QuestionOption{
					{TextOption: "fgh"}, {TextOption: "ghi"},
					{TextOption: "hij", Answer: true}, {TextOption: "ijk"},
				},
			},
		},
	}

	populateDatabase(t, instance.database, getCreator(userID))

	// Act
	response, err := performRequest(http.MethodPost, ts.URL, "api/v1/quizzes", token, input)

	// Assert
	assert.NoError(t, err)
	if !assert.NotNil(t, response) {
		t.FailNow()
	}

	assert.Equal(t, http.StatusOK, response.StatusCode)

	var result *domain.Quiz
	if err := instance.database.Preload("MultipleChoiceQuestions.Options").Find(&result).Error; err != nil {
		t.Fatal(err.Error())
	}

	if assert.NotEmpty(t, result) {
		assert.Equal(t, input.Name, result.Name)

		if assert.Len(t, result.MultipleChoiceQuestions, 1) {
			assert.Equal(t, input.MultipleChoiceQuestions[0].Title, result.MultipleChoiceQuestions[0].Title)
			assert.Equal(t, input.MultipleChoiceQuestions[0].Options[0].TextOption, result.MultipleChoiceQuestions[0].Options[0].TextOption)
		}
	}
}

func TestNewServer_PutQuiz_SavesNewQuiz(t *testing.T) {
	// Arrange
	instance := &Server{jwtSecret: "abc", oAuthConfig: &oauth2.Config{ClientID: "abc", ClientSecret: "abc", RedirectURL: "abc"}}
	instance.database = gormtestutil.NewMemoryDatabase(t, gormtestutil.WithName(t.Name()))

	// Test http server
	engine := gin.Default()
	_ = instance.Configure(engine)
	ts := httptest.NewServer(engine)

	userID := uuid.MustParse("7d87bab0-cf2d-45ae-bced-1de22db21a77")
	token, _ := instance.jwtService.GenerateToken(userID.String())

	// Close it in the end
	defer ts.Close()

	input := &inputs.Quiz{
		Name:        "abc",
		Description: "bcd",
		MultipleChoiceQuestions: []*inputs.MultipleChoiceQuestion{
			{
				Title:             "cde",
				Description:       "def",
				DurationInSeconds: 15,
				Category:          "egh",
				Order:             0,
				Options: []*inputs.QuestionOption{
					{TextOption: "fgh"}, {TextOption: "ghi"},
					{TextOption: "hij", Answer: true}, {TextOption: "ijk"},
				},
			},
		},
	}

	populateDatabase(t, instance.database, getCreator(userID))

	// Act
	response, err := performRequest(http.MethodPut, ts.URL, "api/v1/quizzes/3660def9-bd13-4c94-b9cd-d449eef82503", token, input)

	// Assert
	assert.NoError(t, err)
	if !assert.NotNil(t, response) {
		t.FailNow()
	}

	assert.Equal(t, http.StatusOK, response.StatusCode)

	var result *domain.Quiz
	if err := instance.database.Preload("MultipleChoiceQuestions.Options").Where("id = ?", "3660def9-bd13-4c94-b9cd-d449eef82503").First(&result).Error; err != nil {
		t.Fatal(err.Error())
	}

	if assert.NotEmpty(t, result) {
		assert.Equal(t, input.Name, result.Name)

		if assert.Len(t, result.MultipleChoiceQuestions, 1) {
			assert.Equal(t, input.MultipleChoiceQuestions[0].Title, result.MultipleChoiceQuestions[0].Title)
			assert.Equal(t, input.MultipleChoiceQuestions[0].Options[0].TextOption, result.MultipleChoiceQuestions[0].Options[0].TextOption)
		}
	}
}

func TestNewServer_PutQuiz_UpdatesExistingQuiz(t *testing.T) {
	// Arrange
	instance := &Server{jwtSecret: "abc", oAuthConfig: &oauth2.Config{ClientID: "abc", ClientSecret: "abc", RedirectURL: "abc"}}
	instance.database = gormtestutil.NewMemoryDatabase(t, gormtestutil.WithName(t.Name()))

	// Test http server
	engine := gin.Default()
	_ = instance.Configure(engine)
	ts := httptest.NewServer(engine)

	userID := uuid.MustParse("7d87bab0-cf2d-45ae-bced-1de22db21a77")
	token, _ := instance.jwtService.GenerateToken(userID.String())

	// Close it in the end
	defer ts.Close()

	old := &domain.Quiz{
		BaseObject:  domain.BaseObject{ID: uuid.MustParse("3660def9-bd13-4c94-b9cd-d449eef82503")},
		Name:        "old",
		Description: "older",
		Creator:     getCreator(userID),
		MultipleChoiceQuestions: []*domain.MultipleChoiceQuestion{
			{
				BaseQuestion: domain.BaseQuestion{
					Title:             "bet",
					Description:       "better",
					DurationInSeconds: 20,
					Category:          "agh",
					Order:             0,
				},
				Options: []*domain.QuestionOption{
					{TextOption: "def"}, {TextOption: "abc"},
					{TextOption: "slecht"}, {TextOption: "old"},
				},
			},
		},
	}
	instance.database.Create(old)

	input := &inputs.Quiz{
		Name:        "abc",
		Description: "bcd",
		MultipleChoiceQuestions: []*inputs.MultipleChoiceQuestion{
			{
				Title:             "cde",
				Description:       "def",
				DurationInSeconds: 15,
				Category:          "egh",
				Order:             0,
				Options: []*inputs.QuestionOption{
					{TextOption: "fgh"}, {TextOption: "ghi", Answer: true},
					{TextOption: "hij"}, {TextOption: "ijk"},
				},
			},
		},
	}

	// Act
	response, err := performRequest(http.MethodPut, ts.URL, "api/v1/quizzes/3660def9-bd13-4c94-b9cd-d449eef82503", token, input)

	// Assert
	assert.NoError(t, err)
	if !assert.NotNil(t, response) {
		t.FailNow()
	}

	assert.Equal(t, http.StatusOK, response.StatusCode)

	var result *domain.Quiz
	if err := instance.database.Preload("MultipleChoiceQuestions.Options").Where("id = ?", old.ID).First(&result).Error; err != nil {
		t.Fatal(err.Error())
	}

	if assert.NotEmpty(t, result) {
		assert.Equal(t, input.Name, result.Name)

		if assert.Len(t, result.MultipleChoiceQuestions, 1) {
			assert.Equal(t, input.MultipleChoiceQuestions[0].Title, result.MultipleChoiceQuestions[0].Title)
			assert.Equal(t, input.MultipleChoiceQuestions[0].Options[0].TextOption, result.MultipleChoiceQuestions[0].Options[0].TextOption)
		}
	}
}

func TestNewServer_PostGame_CreatesGame(t *testing.T) {
	// Arrange
	instance := &Server{jwtSecret: "abc", oAuthConfig: &oauth2.Config{ClientID: "abc", ClientSecret: "abc", RedirectURL: "abc"}}
	instance.database = gormtestutil.NewMemoryDatabase(t, gormtestutil.WithName(t.Name()))

	// Test http server
	engine := gin.Default()
	_ = instance.Configure(engine)
	ts := httptest.NewServer(engine)

	userID := uuid.MustParse("7d87bab0-cf2d-45ae-bced-1de22db21a77")
	token, _ := instance.jwtService.GenerateToken(userID.String())

	quizzes := []*domain.Quiz{
		{
			BaseObject: domain.BaseObject{ID: uuid.MustParse("25e48148-3225-4ae9-a737-345b099bca72")},
			Name:       "def",
			Creator:    getCreator(userID),
		},
	}

	// Populate database
	populateDatabase(t, instance.database, quizzes...)

	// Close it in the end
	defer ts.Close()

	input := &inputs.Game{PlayerLimit: 20}

	// Act
	response, err := performRequest(http.MethodPost, ts.URL, "api/v1/quizzes/25e48148-3225-4ae9-a737-345b099bca72/games", token, input)

	// Assert
	assert.NoError(t, err)
	if !assert.NotNil(t, response) {
		t.FailNow()
	}

	assert.Equal(t, http.StatusOK, response.StatusCode)

	var result *domain.Game
	ok := gintestutil.Response(t, &result, http.StatusOK, response)

	if ok {
		assert.Equal(t, input.PlayerLimit, result.PlayerLimit)
	}
}

func TestNewServer_StartGame_StartsGame(t *testing.T) {
	// Arrange
	instance := &Server{jwtSecret: "abc", oAuthConfig: &oauth2.Config{ClientID: "abc", ClientSecret: "abc", RedirectURL: "abc"}}
	instance.database = gormtestutil.NewMemoryDatabase(t, gormtestutil.WithName(t.Name()))

	// Test http server
	engine := gin.Default()
	_ = instance.Configure(engine)
	ts := httptest.NewServer(engine)

	userID := uuid.MustParse("7d87bab0-cf2d-45ae-bced-1de22db21a77")
	token, _ := instance.jwtService.GenerateToken(userID.String())

	game := &domain.Game{
		BaseObject: domain.BaseObject{ID: uuid.MustParse("342855cd-332c-4344-955e-a0e63be17f3a")},
		Quiz: &domain.Quiz{
			BaseObject:              domain.BaseObject{ID: uuid.MustParse("25e48148-3225-4ae9-a737-345b099bca72")},
			Name:                    "def",
			Creator:                 getCreator(userID),
			MultipleChoiceQuestions: []*domain.MultipleChoiceQuestion{{}, {}},
		},
	}

	// Populate database
	populateDatabase(t, instance.database, game)

	// Close it in the end
	defer ts.Close()

	// Act
	response, err := performRequest(http.MethodPatch, ts.URL, "api/v1/games/342855cd-332c-4344-955e-a0e63be17f3a?action=start", token, nil)

	// Assert
	assert.NoError(t, err)
	if !assert.NotNil(t, response) {
		t.FailNow()
	}

	assert.Equal(t, http.StatusOK, response.StatusCode)

	var result *domain.Game
	ok := gintestutil.Response(t, &result, http.StatusOK, response)

	if ok {
		assert.Equal(t, game.PlayerLimit, result.PlayerLimit)
		assert.NotEmpty(t, result.Code)
	}
}

func TestNewServer_DeleteGame_DeletesGame(t *testing.T) {
	// Arrange
	instance := &Server{jwtSecret: "abc", oAuthConfig: &oauth2.Config{ClientID: "abc", ClientSecret: "abc", RedirectURL: "abc"}}
	instance.database = gormtestutil.NewMemoryDatabase(t, gormtestutil.WithName(t.Name()))

	// Test http server
	engine := gin.Default()
	_ = instance.Configure(engine)
	ts := httptest.NewServer(engine)

	userID := uuid.MustParse("7d87bab0-cf2d-45ae-bced-1de22db21a77")
	token, _ := instance.jwtService.GenerateToken(userID.String())

	game := &domain.Game{
		BaseObject: domain.BaseObject{ID: uuid.MustParse("342855cd-332c-4344-955e-a0e63be17f3a")},
		StartTime:  time.Now(),
		FinishTime: time.Now(),
		Code:       "AE52DE",
		Quiz: &domain.Quiz{
			BaseObject: domain.BaseObject{ID: uuid.MustParse("25e48148-3225-4ae9-a737-345b099bca72")},
			Name:       "def",
			Creator:    getCreator(userID),
		},
	}

	// Populate database
	populateDatabase(t, instance.database, game)

	// Close it in the end
	defer ts.Close()

	// Act
	response, err := performRequest(http.MethodDelete, ts.URL, "api/v1/games/342855cd-332c-4344-955e-a0e63be17f3a", token, nil)

	// Assert
	assert.NoError(t, err)
	if !assert.NotNil(t, response) {
		t.FailNow()
	}

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.ErrorContains(t, instance.database.First(&domain.Game{}).Error, "not found")
}

func TestNewServer_GetPlayers_ReturnsData(t *testing.T) {
	// Arrange
	instance := &Server{jwtSecret: "abc", oAuthConfig: &oauth2.Config{ClientID: "abc", ClientSecret: "abc", RedirectURL: "abc"}}
	instance.database = gormtestutil.NewMemoryDatabase(t, gormtestutil.WithName(t.Name()))

	// Test http server
	engine := gin.Default()
	_ = instance.Configure(engine)
	ts := httptest.NewServer(engine)

	userID := uuid.MustParse("7d87bab0-cf2d-45ae-bced-1de22db21a77")
	token, _ := instance.jwtService.GenerateToken(userID.String())

	quizzes := []*domain.Quiz{
		{
			BaseObject: domain.BaseObject{ID: uuid.MustParse("25e48148-3225-4ae9-a737-345b099bca72")},
			Name:       "def",
			Creator:    getCreator(userID),
			Games: []*domain.Game{{
				BaseObject: domain.BaseObject{ID: uuid.MustParse("c37077d7-9922-4bea-af99-1968bfec65e0")},
				Code:       "abc",
				Players:    []*domain.Player{{Nickname: "A"}, {Nickname: "B"}, {Nickname: "C"}},
			}},
		},
	}

	// Populate database
	populateDatabase(t, instance.database, quizzes...)

	// Close it in the end
	defer ts.Close()

	// Act
	response, err := performRequest(http.MethodGet, ts.URL, "api/v1/games/c37077d7-9922-4bea-af99-1968bfec65e0/players", token, nil)

	// Assert
	assert.NoError(t, err)
	if !assert.NotNil(t, response) {
		t.FailNow()
	}

	assert.Equal(t, http.StatusOK, response.StatusCode)

	var result []*domain.Player
	ok := gintestutil.Response(t, &result, http.StatusOK, response)

	if ok && assert.Len(t, result, 3) {
		assert.Equal(t, quizzes[0].Games[0].Players[0].Nickname, result[0].Nickname)
		assert.Equal(t, quizzes[0].Games[0].Players[1].Nickname, result[1].Nickname)
		assert.Equal(t, quizzes[0].Games[0].Players[2].Nickname, result[2].Nickname)
	}
}

func TestNewServer_PostPlayer_AddsPlayerToGame(t *testing.T) {
	// Arrange
	instance := &Server{jwtSecret: "abc", oAuthConfig: &oauth2.Config{ClientID: "abc", ClientSecret: "abc", RedirectURL: "abc"}}
	instance.database = gormtestutil.NewMemoryDatabase(t, gormtestutil.WithName(t.Name()))

	// Test http server
	engine := gin.Default()
	_ = instance.Configure(engine)
	ts := httptest.NewServer(engine)

	userID := uuid.MustParse("7d87bab0-cf2d-45ae-bced-1de22db21a77")

	quizzes := []*domain.Quiz{
		{
			BaseObject: domain.BaseObject{ID: uuid.MustParse("25e48148-3225-4ae9-a737-345b099bca72")},
			Name:       "def",
			Creator:    getCreator(userID),
			Games: []*domain.Game{{
				BaseObject:  domain.BaseObject{ID: uuid.MustParse("c37077d7-9922-4bea-af99-1968bfec65e0")},
				Code:        "abc",
				PlayerLimit: 20,
				StartTime:   time.Now(),
				Players:     []*domain.Player{{Nickname: "A"}, {Nickname: "B"}, {Nickname: "C"}},
			}},
		},
	}

	// Populate database
	populateDatabase(t, instance.database, quizzes...)

	// Close it in the end
	defer ts.Close()

	// Act
	response, err := performRequest(http.MethodPost, ts.URL, "api/v1/games/c37077d7-9922-4bea-af99-1968bfec65e0/players", "", nil)

	// Assert
	assert.NoError(t, err)
	if !assert.NotNil(t, response) {
		t.FailNow()
	}

	assert.Equal(t, http.StatusOK, response.StatusCode)

	var result *domain.Player
	ok := gintestutil.Response(t, &result, http.StatusOK, response)
	if ok {
		assert.NotEmpty(t, result.Nickname)
	}
}

func TestNewServer_DeletePlayer_RemovesPlayerFromGame(t *testing.T) {
	// Arrange
	instance := &Server{jwtSecret: "abc", oAuthConfig: &oauth2.Config{ClientID: "abc", ClientSecret: "abc", RedirectURL: "abc"}}
	instance.database = gormtestutil.NewMemoryDatabase(t, gormtestutil.WithName(t.Name()))

	// Test http server
	engine := gin.Default()
	_ = instance.Configure(engine)
	ts := httptest.NewServer(engine)

	userID := uuid.MustParse("7d87bab0-cf2d-45ae-bced-1de22db21a77")

	quizzes := []*domain.Quiz{
		{
			BaseObject: domain.BaseObject{ID: uuid.MustParse("25e48148-3225-4ae9-a737-345b099bca72")},
			Name:       "def",
			Creator:    getCreator(userID),
			Games: []*domain.Game{{
				BaseObject:  domain.BaseObject{ID: uuid.MustParse("c37077d7-9922-4bea-af99-1968bfec65e0")},
				Code:        "abc",
				PlayerLimit: 20,
				StartTime:   time.Now(),
				Players: []*domain.Player{
					{BaseObject: domain.BaseObject{ID: uuid.MustParse("c23330d9-3d58-45cd-a49e-8085f4c15439")}, Nickname: "A"},
				},
			}},
		},
	}

	// Populate database
	populateDatabase(t, instance.database, quizzes...)

	// Close it in the end
	defer ts.Close()

	// Act
	response, err := performRequest(http.MethodDelete, ts.URL, "api/v1/players/c23330d9-3d58-45cd-a49e-8085f4c15439", "", nil)

	// Assert
	assert.NoError(t, err)
	if !assert.NotNil(t, response) {
		t.FailNow()
	}

	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestNewServer_GetPublicQuiz_ReturnsExpectedQuiz(t *testing.T) {
	// Arrange
	instance := &Server{jwtSecret: "abc", oAuthConfig: &oauth2.Config{ClientID: "abc", ClientSecret: "abc", RedirectURL: "abc"}}
	instance.database = gormtestutil.NewMemoryDatabase(t, gormtestutil.WithName(t.Name()))

	// Test http server
	engine := gin.Default()
	_ = instance.Configure(engine)
	ts := httptest.NewServer(engine)

	userID := uuid.MustParse("7d87bab0-cf2d-45ae-bced-1de22db21a77")

	quizzes := []*domain.Quiz{
		{
			BaseObject: domain.BaseObject{ID: uuid.MustParse("25e48148-3225-4ae9-a737-345b099bca72")},
			Name:       "def",
			Creator:    getCreator(userID),
			Games: []*domain.Game{{
				BaseObject:  domain.BaseObject{ID: uuid.MustParse("c37077d7-9922-4bea-af99-1968bfec65e0")},
				Code:        "abc",
				PlayerLimit: 20,
				StartTime:   time.Now(),
				Players: []*domain.Player{
					{BaseObject: domain.BaseObject{ID: uuid.MustParse("c23330d9-3d58-45cd-a49e-8085f4c15439")}, Nickname: "A"},
				},
			}},
		},
	}

	// Populate database
	populateDatabase(t, instance.database, quizzes...)

	// Close it in the end
	defer ts.Close()

	// Act
	response, err := performRequest(http.MethodGet, ts.URL, "api/v1/games/c37077d7-9922-4bea-af99-1968bfec65e0/quiz", "", nil)

	// Assert
	assert.NoError(t, err)
	if !assert.NotNil(t, response) {
		t.FailNow()
	}

	assert.Equal(t, http.StatusOK, response.StatusCode)

	var result *outputs.OutputQuiz
	ok := gintestutil.Response(t, &result, http.StatusOK, response)
	if ok {
		assert.Equal(t, quizzes[0].Name, result.Name)
		assert.Equal(t, quizzes[0].ID, result.ID)
		assert.Equal(t, quizzes[0].Description, result.Description)
	}
}

// This tests:
// - Create game
// - Start game
// - 2 players join
// - Next question
// - 2 answers
// - Nest Question
// - 2 answers
// - Finish
func TestNewServer_GameFlow_Works(t *testing.T) {
	// Arrange
	instance := &Server{jwtSecret: "abc", oAuthConfig: &oauth2.Config{ClientID: "abc", ClientSecret: "abc", RedirectURL: "abc"}}
	instance.database = gormtestutil.NewMemoryDatabase(t, gormtestutil.WithName(t.Name()), gormtestutil.WithSingularConnection())

	// Test http server
	engine := gin.Default()
	_ = instance.Configure(engine)
	ts := httptest.NewServer(engine)

	userID := uuid.MustParse("7d87bab0-cf2d-45ae-bced-1de22db21a77")
	token, _ := instance.jwtService.GenerateToken(userID.String())

	quiz := &domain.Quiz{
		BaseObject: domain.BaseObject{ID: uuid.MustParse("25e48148-3225-4ae9-a737-345b099bca72")},
		Name:       "def",
		Creator:    getCreator(userID),
		MultipleChoiceQuestions: []*domain.MultipleChoiceQuestion{
			{
				BaseQuestion: domain.BaseQuestion{
					BaseObject:        domain.BaseObject{ID: uuid.MustParse("5413ddc1-986c-43cf-8150-3aa3eb1e5f4f")},
					Order:             0,
					DurationInSeconds: 4,
				},
				Options: []*domain.QuestionOption{
					{BaseObject: domain.BaseObject{ID: uuid.MustParse("4f95d9ce-a608-4292-b3f6-18b4b7939135")}},
				},
			},
			{
				BaseQuestion: domain.BaseQuestion{
					BaseObject:        domain.BaseObject{ID: uuid.MustParse("c847e53b-9dd6-4636-99be-6cf18243d598")},
					Order:             1,
					DurationInSeconds: 4,
				},
				Options: []*domain.QuestionOption{
					{BaseObject: domain.BaseObject{ID: uuid.MustParse("7b7a4cdd-622a-4a57-adb4-064ada2bc4fa")}},
				},
			},
		},
	}

	// Populate database
	populateDatabase(t, instance.database, quiz)

	// Close it in the end
	defer ts.Close()

	// Act
	gameRes, err := performRequest(http.MethodPost, ts.URL, "api/v1/quizzes/25e48148-3225-4ae9-a737-345b099bca72/games", token, inputs.Game{PlayerLimit: 5})
	gameID := getValue(t, gameRes, err, func(t domain.Game) uuid.UUID {
		return t.ID
	})
	_, _ = performRequest(http.MethodPatch, ts.URL, fmt.Sprintf("api/v1/games/%s?action=start", gameID), token, nil)

	creatorUrl := fmt.Sprintf("ws%s/api/v1/games/%s/connection", strings.TrimPrefix(ts.URL, "http"), gameID)
	creatorSocket, _, creatorErr := websocket.DefaultDialer.Dial(creatorUrl, http.Header{"Authorization": []string{"Bearer " + token}})

	// Player 1 connection
	player1Res, _ := performRequest(http.MethodPost, ts.URL, fmt.Sprintf("api/v1/games/%s/players", gameID), "", nil)
	player1ID := getValue(t, player1Res, err, func(t domain.Player) uuid.UUID {
		return t.ID
	})
	player1Url := fmt.Sprintf("ws%s/api/v1/games/%s/players/%s/connection", strings.TrimPrefix(ts.URL, "http"), gameID, player1ID)
	player1Socket, _, p1Err := websocket.DefaultDialer.Dial(player1Url, nil)

	// Player 2 connection
	player2Res, _ := performRequest(http.MethodPost, ts.URL, fmt.Sprintf("api/v1/games/%s/players", gameID), "", nil)
	player2ID := getValue(t, player2Res, err, func(t domain.Player) uuid.UUID {
		return t.ID
	})
	player2Url := fmt.Sprintf("ws%s/api/v1/games/%s/players/%s/connection", strings.TrimPrefix(ts.URL, "http"), gameID, player2ID)
	player2Socket, _, p2Err := websocket.DefaultDialer.Dial(player2Url, nil)

	if p1Err != nil || p2Err != nil || creatorErr != nil {
		t.Fatal(p1Err, p2Err, creatorErr)
	}
	defer creatorSocket.Close()
	defer player1Socket.Close()
	defer player2Socket.Close()

	// Play game
	_ = creatorSocket.WriteJSON(&coordinator.CreatorMessage{Action: coordinator.NextQuestionAction})

	// Wait a second to propagate
	time.Sleep(500 * time.Millisecond)

	// Answer questions
	_ = player1Socket.WriteJSON(&coordinator.PlayerMessage{Action: coordinator.AnswerAction, Content: json.RawMessage("{\"optionID\": \"4f95d9ce-a608-4292-b3f6-18b4b7939135\"}")})
	_ = player2Socket.WriteJSON(&coordinator.PlayerMessage{Action: coordinator.AnswerAction, Content: json.RawMessage("{\"optionID\": \"4f95d9ce-a608-4292-b3f6-18b4b7939135\"}")})

	// Deadline
	time.Sleep(4 * time.Second)

	// Next question
	_ = creatorSocket.WriteJSON(&coordinator.CreatorMessage{Action: coordinator.NextQuestionAction})

	// Wait a second to propagate
	time.Sleep(500 * time.Millisecond)

	// Answer questions
	_ = player1Socket.WriteJSON(&coordinator.PlayerMessage{Action: coordinator.AnswerAction, Content: json.RawMessage("{\"optionID\": \"7b7a4cdd-622a-4a57-adb4-064ada2bc4fa\"}")})
	_ = player2Socket.WriteJSON(&coordinator.PlayerMessage{Action: coordinator.AnswerAction, Content: json.RawMessage("{\"optionID\": \"7b7a4cdd-622a-4a57-adb4-064ada2bc4fa\"}")})

	// Deadline
	time.Sleep(4 * time.Second)

	// Finish the game
	_ = creatorSocket.WriteJSON(&coordinator.CreatorMessage{Action: coordinator.FinishGameAction})

	// Wait a second to propagate
	time.Sleep(1 * time.Second)

	// Verify
	res, err := performRequest(http.MethodGet, ts.URL, fmt.Sprintf("api/v1/games/%s", gameID), token, nil)

	// Assert
	assert.NoError(t, err)

	// Check final result
	var result *domain.Game
	ok := gintestutil.Response(t, &result, http.StatusOK, res)

	if ok {
		assert.False(t, result.FinishTime.IsZero())
		assert.Len(t, result.Players, 2)
		assert.Len(t, result.Answers, 4)
	}
}
