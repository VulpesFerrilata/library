package app_error

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/kataras/iris/v12"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/go-playground/validator.v9"
)

var (
	ErrInvalidValidationType error = errors.New("invalid validation type")
)

func NewInputValidationError(s interface{}, fieldErrors validator.ValidationErrors) AppError {
	return newValidationError(inputValidation, s, fieldErrors)
}

func NewEntityValidationError(s interface{}, fieldErrors validator.ValidationErrors) AppError {
	return newValidationError(entityValidation, s, fieldErrors)
}

func newValidationError(validationType validationType, s interface{}, fieldErrors validator.ValidationErrors) AppError {
	var namespace string
	if t := reflect.TypeOf(s); t.Kind() == reflect.Ptr {
		namespace = t.Elem().Name()
	} else {
		namespace = t.Name()
	}
	return &validationError{
		validationType: validationType,
		namespace:      namespace,
		fieldErrors:    fieldErrors,
	}
}

type validationType int

const (
	inputValidation validationType = iota
	entityValidation
)

type validationError struct {
	validationType validationType
	namespace      string
	fieldErrors    validator.ValidationErrors
}

func (ve validationError) Error() string {
	builder := new(strings.Builder)

	builder.WriteString(fmt.Sprintf("one or more fields in %s contain invalid data", ve.namespace))
	builder.WriteString("\n")
	builder.WriteString(ve.fieldErrors.Error())

	return strings.TrimSpace(builder.String())
}

func (ve validationError) Problem(trans ut.Translator) (iris.Problem, error) {
	problem := iris.NewProblem()
	problem.Type("about:blank")

	switch ve.validationType {
	case inputValidation:
		problem.Status(iris.StatusUnprocessableEntity)
		title, err := trans.T("unprocessable-entity-error")
		if err != nil {
			return nil, err
		}
		problem.Title(title)
	case entityValidation:
		problem.Status(iris.StatusInternalServerError)
		title, err := trans.T("internal-server-error")
		if err != nil {
			return nil, err
		}
		problem.Title(title)
	default:
		return nil, ErrInvalidValidationType
	}

	detail, err := trans.T("validation-error-detail", ve.namespace)
	if err != nil {
		return nil, err
	}
	problem.Detail(detail)

	var errors []string
	for _, fieldError := range ve.fieldErrors {
		fieldErrorTrans := fieldError.Translate(trans)
		errors = append(errors, fieldErrorTrans)
	}
	problem.Key("errors", errors)

	return problem, nil
}

func (ve validationError) Status(trans ut.Translator) (*status.Status, error) {
	detail, err := trans.T("validation-error-detail", ve.namespace)
	if err != nil {
		return nil, err
	}
	var code codes.Code
	switch ve.validationType {
	case inputValidation:
		code = codes.InvalidArgument
	case entityValidation:
		code = codes.Internal
	default:
		return nil, ErrInvalidValidationType
	}
	stt := status.New(code, detail)
	badRequest := &errdetails.BadRequest{}
	fieldErrors := ve.fieldErrors.Translate(trans)
	for namespace, fieldError := range fieldErrors {
		fieldViolation := &errdetails.BadRequest_FieldViolation{
			Field:       namespace,
			Description: fieldError,
		}

		badRequest.FieldViolations = append(badRequest.FieldViolations, fieldViolation)
	}

	return stt.WithDetails(badRequest)
}

func (ve validationError) Message(trans ut.Translator) (string, error) {
	builder := new(strings.Builder)

	for _, fieldError := range ve.fieldErrors {
		builder.WriteString(fieldError.Translate(trans))
		builder.WriteString("\n")
	}

	return builder.String(), nil
}
