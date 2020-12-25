package app_errors

import ut "github.com/go-playground/universal-translator"

type BusinessRuleError interface {
	error
	Translate(trans ut.Translator) (string, error)
}
