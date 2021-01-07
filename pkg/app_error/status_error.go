package app_error

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/kataras/iris/v12"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StatusError interface {
	GRPCStatus() *status.Status
}

func NewStatusError(stt *status.Status) WebError {
	return &statusError{
		stt: stt,
	}
}

type statusError struct {
	stt *status.Status
}

func (se *statusError) Error() string {
	return se.stt.Err().Error()
}

func (se *statusError) Problem(trans ut.Translator) (iris.Problem, error) {
	problem := iris.NewProblem()
	problem.Type("about:blank")
	switch se.stt.Code() {
	case codes.InvalidArgument:
		problem.Status(iris.StatusUnprocessableEntity)
	case codes.NotFound:
		problem.Status(iris.StatusUnprocessableEntity)
	default:
		problem.Status(iris.StatusInternalServerError)
	}
	problem.Detail(se.stt.Message())

	errs := make([]string, 0)
	for _, detail := range se.stt.Details() {
		switch detailType := detail.(type) {
		case errdetails.BadRequest:
			for _, fieldViolation := range detailType.GetFieldViolations() {
				errs = append(errs, fieldViolation.GetDescription())
			}
		case errdetails.PreconditionFailure:
			for _, violation := range detailType.GetViolations() {
				errs = append(errs, violation.GetDescription())
			}
		}
	}
	if len(errs) > 0 {
		problem.Key("errors", errs)
	}

	return problem, nil
}
