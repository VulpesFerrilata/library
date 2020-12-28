package app_error

import (
	"fmt"

	ut "github.com/go-playground/universal-translator"
)

func NewAlreadyExistsError(name string) BusinessRuleError {
	return &AlreadyExistsError{
		name: name,
	}
}

type AlreadyExistsError struct {
	name string
}

func (aee AlreadyExistsError) Error() string {
	return fmt.Sprintf("%s is already exists", aee.name)
}

func (aee AlreadyExistsError) Translate(trans ut.Translator) (string, error) {
	return trans.T("already-exists-error", aee.name)
}
