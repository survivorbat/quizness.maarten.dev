package inputs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMultipleChoiceQuestion_HasOneAnswer(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		input    *MultipleChoiceQuestion
		expected bool
	}{
		"no answers": {
			input: &MultipleChoiceQuestion{
				Options: []*QuestionOption{{}, {}},
			},
		},
		"1 answer": {
			input: &MultipleChoiceQuestion{
				Options: []*QuestionOption{{Answer: true}, {}},
			},
			expected: true,
		},
		"multiple answer": {
			input: &MultipleChoiceQuestion{
				Options: []*QuestionOption{{Answer: true}, {Answer: true}},
			},
		},
	}

	for name, testData := range tests {
		testData := testData
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			// Act
			result := testData.input.hasOneAnswer()

			// Assert
			assert.Equal(t, testData.expected, result)
		})
	}
}

func TestQuiz_HasValidOrder_ReturnsExpectedResult(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input    *Quiz
		expected bool
	}{
		"none": {
			input:    &Quiz{},
			expected: true,
		},
		"3 questions": {
			input: &Quiz{
				MultipleChoiceQuestions: []*MultipleChoiceQuestion{
					{Order: 0},
					{Order: 2},
					{Order: 1},
				},
			},
			expected: true,
		},
		"5 invalid questions": {
			input: &Quiz{
				MultipleChoiceQuestions: []*MultipleChoiceQuestion{
					{Order: 0},
					{Order: 1},
					{Order: 2},
					{Order: 3},
					{Order: 5},
				},
			},
		},
		"complete nonsense": {
			input: &Quiz{
				MultipleChoiceQuestions: []*MultipleChoiceQuestion{
					{Order: 0},
					{Order: 5},
					{Order: 21},
					{Order: 8},
					{Order: 70},
				},
			},
		},
	}

	for name, testData := range tests {
		testData := testData
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			// Act
			result := testData.input.hasValidOrder()

			// Assert
			assert.Equal(t, testData.expected, result)
		})
	}
}
