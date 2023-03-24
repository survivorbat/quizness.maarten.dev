package inputs

import (
	"github.com/go-playground/validator/v10"
	"reflect"
)

type ReportedError struct {
	field           any
	fieldName       string
	structFieldName string
	tag             string
	param           string
}

type StructLevelMock struct {
	validator.StructLevel

	currentValue   reflect.Value
	reportedErrors []ReportedError
}

func (s *StructLevelMock) Validator() *validator.Validate {
	return validator.New()
}

func (s *StructLevelMock) Current() reflect.Value {
	return s.currentValue
}

func (s *StructLevelMock) ReportError(field any, fieldName, structFieldName string, tag, param string) {
	s.reportedErrors = append(s.reportedErrors, ReportedError{field, fieldName, structFieldName, tag, param})
}
