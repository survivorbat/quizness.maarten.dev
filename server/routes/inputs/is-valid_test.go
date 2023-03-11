package inputs

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type testInput struct {
	returnOk          bool
	returnField       any
	returnFieldName   string
	returnStructField string
	returnTag         string
	returnParam       string
}

func (t *testInput) IsValid() (bool, any, string, string, string, string) {
	return t.returnOk, t.returnField, t.returnFieldName, t.returnStructField, t.returnTag, t.returnParam
}

func TestIsValidator_ReportsAnyError(t *testing.T) {
	t.Parallel()
	// Arrange
	input := &StructLevelMock{
		currentValue: reflect.ValueOf(&testInput{
			returnOk:          true,
			returnField:       "a",
			returnParam:       "b",
			returnTag:         "c",
			returnStructField: "d",
			returnFieldName:   "e",
		}),
	}

	// Act
	IsValidator(input)

	// Assert
	expected := []ReportedError{
		{
			field:           "a",
			fieldName:       "e",
			structFieldName: "d",
			tag:             "c",
			param:           "b",
		},
	}
	assert.Equal(t, expected, input.reportedErrors)
}
