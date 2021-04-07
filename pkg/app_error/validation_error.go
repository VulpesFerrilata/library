package app_error

import (
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/kataras/iris/v12"
	"github.com/pkg/errors"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/go-playground/validator.v9"
)

func NewValidationError(fieldErrors validator.ValidationErrors) AppError {
	return &validationError{
		fieldErrors: fieldErrors,
	}
}

type validationError struct {
	fieldErrors validator.ValidationErrors
}

func (ve validationError) Error() string {
	builder := new(strings.Builder)

	builder.WriteString("one or more fields contain invalid data")
	builder.WriteString("\n")
	builder.WriteString(ve.fieldErrors.Error())

	return builder.String()
}

func (ve validationError) Problem(trans ut.Translator) (iris.Problem, error) {
	problem := iris.NewProblem()
	problem.Type("about:blank")

	problem.Status(iris.StatusUnprocessableEntity)
	detail, err := trans.T("validation-error")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	problem.Detail(detail)

	errors := make([]string, 0)
	for _, fieldError := range ve.fieldErrors {
		fieldErrorTrans := fieldError.Translate(trans)
		errors = append(errors, fieldErrorTrans)
	}
	problem.Key("errors", errors)

	return problem, nil
}

func (ve validationError) Status(trans ut.Translator) (*status.Status, error) {
	detail, err := trans.T("validation-error")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	stt := status.New(codes.InvalidArgument, detail)

	badRequest := &errdetails.BadRequest{}
	for _, fieldError := range ve.fieldErrors {
		fieldViolation := &errdetails.BadRequest_FieldViolation{
			Field:       fieldError.Field(),
			Description: fieldError.Translate(trans),
		}

		badRequest.FieldViolations = append(badRequest.FieldViolations, fieldViolation)
	}

	return stt.WithDetails(badRequest)
}

func (ve validationError) Message(trans ut.Translator) (string, error) {
	builder := new(strings.Builder)

	detail, err := trans.T("validation-error")
	if err != nil {
		return "", errors.WithStack(err)
	}
	builder.WriteString(detail)
	for _, fieldError := range ve.fieldErrors {
		builder.WriteString("\n")
		builder.WriteString(fieldError.Translate(trans))
	}

	return builder.String(), nil
}
