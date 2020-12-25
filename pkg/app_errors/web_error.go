package app_errors

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/kataras/iris/v12"
)

type WebError interface {
	error
	Problem(trans ut.Translator) (iris.Problem, error)
}
