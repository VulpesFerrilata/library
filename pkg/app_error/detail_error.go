package app_error

import (
	ut "github.com/go-playground/universal-translator"
)

type DetailError interface {
	error
	Translate(trans ut.Translator) (string, error)
}
