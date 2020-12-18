package app_errors

import (
	"bytes"
	"errors"
	"fmt"
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

func NewValidationError(validationType validationType, name string, fieldErrors validator.ValidationErrors) AppError {
	return &validationError{
		validationType: validationType,
		name:           name,
		fieldErrors:    fieldErrors,
	}
}

type validationType int

const (
	InputValidation validationType = iota
	EntityValidation
)

func (vt validationType) String() string {
	return [...]string{"request", "entity"}[vt]
}

type validationError struct {
	validationType validationType
	name           string
	fieldErrors    validator.ValidationErrors
}

func (ve *validationError) Error() string {
	buff := bytes.NewBufferString("")

	buff.WriteString(fmt.Sprintf("one or more fields in %s %s contain invalid data", ve.name, ve.validationType.String()))
	buff.WriteString(ve.fieldErrors.Error())

	return strings.TrimSpace(buff.String())
}

func (ve *validationError) Problem(trans ut.Translator) (iris.Problem, error) {
	problem := iris.NewProblem()
	problem.Type("about:blank")

	switch ve.validationType {
	case InputValidation:
		problem.Status(iris.StatusUnprocessableEntity)
		title, err := trans.T("unprocessable-entity-error")
		if err != nil {
			return nil, err
		}
		problem.Title(title)
	case EntityValidation:
		problem.Status(iris.StatusInternalServerError)
		title, err := trans.T("internal-server-error")
		if err != nil {
			return nil, err
		}
		problem.Title(title)
	default:
		return nil, ErrInvalidValidationType
	}
	detail, _ := trans.T("validation-error-detail", ve.name, ve.validationType.String())
	problem.Detail(detail)
	fieldErrors := ve.fieldErrors.Translate(trans)
	problem.Key("errors", fieldErrors)
	return problem, nil
}

func (ve *validationError) Status(trans ut.Translator) (*status.Status, error) {
	detail, _ := trans.T("validation-error-detail", ve.name, ve.validationType.String())
	var code codes.Code
	switch ve.validationType {
	case InputValidation:
		code = codes.InvalidArgument
	case EntityValidation:
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
