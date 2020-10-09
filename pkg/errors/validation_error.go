package errors

import (
	"bytes"
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/kataras/iris/v12"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ValidationError interface {
	Error
	WithFieldError(fieldErr string)
	HasErrors() bool
}

func NewValidationError() ValidationError {
	return &validationError{}
}

type validationError struct {
	fieldErrs []string
}

func (ve validationError) Error() string {
	buff := bytes.NewBufferString("")

	for _, fieldErr := range ve.fieldErrs {
		buff.WriteString(fieldErr)
		buff.WriteString("\n")
	}

	return strings.TrimSpace(buff.String())
}

func (ve *validationError) WithFieldError(fieldErr string) {
	ve.fieldErrs = append(ve.fieldErrs, fieldErr)
}

func (ve validationError) HasErrors() bool {
	if len(ve.fieldErrs) > 0 {
		return true
	}
	return false
}

func (ve validationError) ToProblem(trans ut.Translator) iris.Problem {
	problem := iris.NewProblem()
	problem.Status(iris.StatusUnprocessableEntity)
	problem.Type("about:blank")
	title, _ := trans.T("validation-error")
	problem.Title(title)
	detail, _ := trans.T("validation-error-detail")
	problem.Detail(detail)
	problem.Key("errors", ve.fieldErrs)
	return problem
}

func (ve validationError) ToStatus(trans ut.Translator) *status.Status {
	title, _ := trans.T("validation-error")
	detail, _ := trans.T("validation-error-detail")
	stt := status.New(codes.FailedPrecondition, detail)
	preconditionFailure := &errdetails.PreconditionFailure{}
	for _, fieldErr := range ve.fieldErrs {
		preconditionViolation := &errdetails.PreconditionFailure_Violation{
			Type:        title,
			Description: fieldErr,
		}
		preconditionFailure.Violations = append(preconditionFailure.Violations, preconditionViolation)
	}
	detailStt, err := stt.WithDetails(preconditionFailure)
	if err != nil {
		return stt
	}
	return detailStt
}
