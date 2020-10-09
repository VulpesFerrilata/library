package errors

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/kataras/iris/v12"
	"google.golang.org/grpc/status"
)

type Error interface {
	error
	ToProblem(trans ut.Translator) iris.Problem
	ToStatus(trans ut.Translator) *status.Status
}
