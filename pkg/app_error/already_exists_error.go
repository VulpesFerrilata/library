package app_error

import (
	"fmt"

	ut "github.com/go-playground/universal-translator"
	"github.com/pkg/errors"
)

func IsAlreadyExistsError(err error) bool {
	_, ok := err.(*alreadyExistsError)
	return ok
}

func NewAlreadyExistsError(name string) BusinessRuleError {
	return &alreadyExistsError{
		name: name,
	}
}

type alreadyExistsError struct {
	name string
}

func (a alreadyExistsError) Error() string {
	return fmt.Sprintf("%s is already exists", a.name)
}

func (a alreadyExistsError) Translate(trans ut.Translator) (string, error) {
	detail, err := trans.T("already-exists-error", a.name)
	return detail, errors.WithStack(err)
}
