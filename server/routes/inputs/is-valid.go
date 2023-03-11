package inputs

import (
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type IsValid interface {
	IsValid() (bool, any, string, string, string, string)
}

func IsValidator(sl validator.StructLevel) {
	isValid, ok := sl.Current().Interface().(IsValid)
	if !ok {
		logrus.Errorf("Not IsValid")
		return
	}

	if ok, field, fieldName, structFieldName, tag, param := isValid.IsValid(); ok {
		sl.ReportError(field, fieldName, structFieldName, tag, param)
	}
}
