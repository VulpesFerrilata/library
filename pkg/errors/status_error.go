package errors

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/kataras/iris/v12"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewStatusError(stt *status.Status) (Error, bool) {
	if stt.Code() == codes.OK {
		return nil, false
	}
	return &StatusError{
		stt: stt,
	}, true
}

type StatusError struct {
	stt *status.Status
}

func (se StatusError) Error() string {
	return se.stt.Err().Error()
}

func (se StatusError) ToProblem(trans ut.Translator) iris.Problem {
	problem := iris.NewProblem()
	problem.Type("about:blank")
	switch se.stt.Code() {
	case codes.InvalidArgument:
		problem.Status(iris.StatusUnprocessableEntity)
		title, _ := trans.T("validation-error")
		problem.Title(title)
		detail, _ := trans.T("validation-error-detail")
		problem.Detail(detail)
	case codes.NotFound:
		problem.Status(iris.StatusNotFound)
		title, _ := trans.T("not-found-error")
		problem.Title(title)
		problem.Detail(se.stt.Message())
	default:
		problem.Status(iris.StatusInternalServerError)
		title, _ := trans.T("internal-error")
		problem.Title(title)
		problem.Detail(se.stt.Message())
	}

	for _, detail := range se.stt.Details() {
		switch detailType := detail.(type) {
		case errdetails.PreconditionFailure:
			errs := make([]string, 0)
			for _, preconditionViolation := range detailType.GetViolations() {
				errs = append(errs, preconditionViolation.Description)
			}
			problem.Key("errors", errs)
		}
	}

	return problem
}

func (se StatusError) ToStatus(trans ut.Translator) *status.Status {
	return se.stt
}
