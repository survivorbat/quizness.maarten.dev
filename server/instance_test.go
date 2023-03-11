package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/survivorbat/qq.maarten.dev/server/domain"
	"github.com/survivorbat/qq.maarten.dev/server/routes/inputs"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// populateDatabase Populates the database with the given data
func populateDatabase[T any](t *testing.T, database *gorm.DB, data ...T) {
	if err := database.Model(new(T)).CreateInBatches(data, 100).Error; err != nil {
		t.Fatal(err.Error())
	}
}

func TestNewServer_GetQuizzes_ReturnsData(t *testing.T) {
	// Arrange
	databaseOpen = sqlite.Open
	connection := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	instance, _ := NewServer(connection, "abc", "abc", "abc", "abc")

	// Test http server
	engine := gin.Default()
	_ = instance.Configure(engine)
	ts := httptest.NewServer(engine)

	userID := uuid.MustParse("7d87bab0-cf2d-45ae-bced-1de22db21a77")
	token, _ := instance.jwtService.GenerateToken(userID.String())

	quizzes := []*domain.Quiz{
		{Name: "def", CreatorID: uuid.MustParse("dc0057c9-553d-40aa-a0bf-6fb98990c634")},
		{
			Name:      "abc",
			CreatorID: userID,
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
	populateDatabase(t, instance.database, quizzes...)

	// Close it in the end
	defer ts.Close()

	requestUrl, _ := url.Parse(fmt.Sprintf("%s/api/v1/quizzes", ts.URL))
	request := &http.Request{
		Method: http.MethodGet,
		Header: map[string][]string{"Authorization": {fmt.Sprintf("Bearer %s", token)}},
		URL:    requestUrl,
	}

	// Act
	response, err := http.DefaultClient.Do(request)

	// Assert
	assert.NoError(t, err)
	if !assert.NotNil(t, response) {
		t.FailNow()
	}

	assert.Equal(t, http.StatusOK, response.StatusCode)

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err.Error())
	}

	var result []*domain.Quiz
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatal(err.Error())
	}

	if assert.Len(t, result, 1) {
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

			// Test http server
			engine := gin.Default()
			_ = instance.Configure(engine)
			ts := httptest.NewServer(engine)

			userID := uuid.MustParse("7d87bab0-cf2d-45ae-bced-1de22db21a77")
			token, _ := instance.jwtService.GenerateToken(userID.String())

			// Close it in the end
			defer ts.Close()

			inputJson, _ := json.Marshal(testData.input)

			requestUrl, _ := url.Parse(fmt.Sprintf("%s/api/v1/quizzes", ts.URL))
			request := &http.Request{
				Method: http.MethodPost,
				Header: map[string][]string{"Authorization": {fmt.Sprintf("Bearer %s", token)}},
				URL:    requestUrl,
				Body:   io.NopCloser(bytes.NewBuffer(inputJson)),
			}

			// Act
			response, err := http.DefaultClient.Do(request)

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
	databaseOpen = sqlite.Open
	connection := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	instance, _ := NewServer(connection, "abc", "abc", "abc", "abc")

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

	inputJson, _ := json.Marshal(input)

	requestUrl, _ := url.Parse(fmt.Sprintf("%s/api/v1/quizzes", ts.URL))
	request := &http.Request{
		Method: http.MethodPost,
		Header: map[string][]string{"Authorization": {fmt.Sprintf("Bearer %s", token)}},
		URL:    requestUrl,
		Body:   io.NopCloser(bytes.NewBuffer(inputJson)),
	}

	// Act
	response, err := http.DefaultClient.Do(request)

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
	databaseOpen = sqlite.Open
	connection := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	instance, _ := NewServer(connection, "abc", "abc", "abc", "abc")

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

	inputJson, _ := json.Marshal(input)

	requestUrl, _ := url.Parse(fmt.Sprintf("%s/api/v1/quizzes/3660def9-bd13-4c94-b9cd-d449eef82503", ts.URL))
	request := &http.Request{
		Method: http.MethodPut,
		Header: map[string][]string{"Authorization": {fmt.Sprintf("Bearer %s", token)}},
		URL:    requestUrl,
		Body:   io.NopCloser(bytes.NewBuffer(inputJson)),
	}

	// Act
	response, err := http.DefaultClient.Do(request)

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
	databaseOpen = sqlite.Open
	connection := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	instance, _ := NewServer(connection, "abc", "abc", "abc", "abc")

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
		Creator:     &domain.Creator{BaseObject: domain.BaseObject{ID: userID}},
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

	inputJson, _ := json.Marshal(input)

	requestUrl, _ := url.Parse(fmt.Sprintf("%s/api/v1/quizzes/3660def9-bd13-4c94-b9cd-d449eef82503", ts.URL))
	request := &http.Request{
		Method: http.MethodPut,
		Header: map[string][]string{"Authorization": {fmt.Sprintf("Bearer %s", token)}},
		URL:    requestUrl,
		Body:   io.NopCloser(bytes.NewBuffer(inputJson)),
	}

	// Act
	response, err := http.DefaultClient.Do(request)

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
