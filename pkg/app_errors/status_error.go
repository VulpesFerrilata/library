package app_errors

import (
	"errors"

	ut "github.com/go-playground/universal-translator"
	"github.com/kataras/iris/v12"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/go-playground/validator.v9"
)

var (
	ErrInvalidStatusCode error = errors.New("invalid status code")
)

func NewStatusError(stt *status.Status) WebError {
	return &statusError{
		stt: stt,
	}
}

type statusError struct {
	stt *status.Status
}

func (se *statusError) Problem(trans ut.Translator) (iris.Problem, error) {
	problem := iris.NewProblem()
	problem.Type("about:blank")
	switch se.stt.Code() {
	case codes.InvalidArgument:
		problem.Status(iris.StatusUnprocessableEntity)
		title, err := trans.T("unprocessable-entity-error")
		if err != nil {
			return nil, err
		}
		problem.Title(title)
	case codes.NotFound:
		problem.Status(iris.StatusNotFound)
		title, err := trans.T("not-found-error")
		if err != nil {
			return nil, err
		}
		problem.Title(title)
	case codes.Internal:
		problem.Status(iris.StatusInternalServerError)
		title, err := trans.T("internal-server-error")
		if err != nil {
			return nil, err
		}
		problem.Title(title)
	default:
		return nil, ErrInvalidStatusCode
	}
	problem.Detail(se.stt.Message())

	for _, detail := range se.stt.Details() {
		switch detailType := detail.(type) {
		case errdetails.BadRequest:
			fieldErrors := make(validator.ValidationErrorsTranslations)
			for _, fieldViolation := range detailType.GetFieldViolations() {
				fieldErrors[fieldViolation.GetField()] = fieldViolation.GetDescription()
			}
			problem.Key("errors", fieldErrors)
		}
	}

	return problem, nil
}
