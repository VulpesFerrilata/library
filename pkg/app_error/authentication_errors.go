package app_error

import (
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/kataras/iris/v12"
	"github.com/pkg/errors"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewAuthenticationErrors() WithDetailAppError {
	return &authenticationErrors{
		detailErrs: make([]DetailError, 0),
	}
}

type authenticationErrors struct {
	detailErrs []DetailError
}

func (a *authenticationErrors) AddDetailError(detailErr DetailError) {
	a.detailErrs = append(a.detailErrs, detailErr)
}

func (a authenticationErrors) Error() string {
	builder := new(strings.Builder)

	builder.WriteString("authentication failed")
	for _, detailErr := range a.detailErrs {
		builder.WriteString("\n")
		builder.WriteString(detailErr.Error())
	}

	return builder.String()
}

func (a authenticationErrors) Problem(trans ut.Translator) (iris.Problem, error) {
	problem := iris.NewProblem()
	problem.Type("about:blank")

	problem.Status(iris.StatusUnprocessableEntity)
	detail, err := trans.T("authentication-error")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	problem.Detail(detail)

	var errs []string
	for _, detailErr := range a.detailErrs {
		authenticationErrTrans, err := detailErr.Translate(trans)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		errs = append(errs, authenticationErrTrans)
	}
	problem.Key("errors", errs)

	return problem, nil
}

func (a authenticationErrors) Status(trans ut.Translator) (*status.Status, error) {
	detail, err := trans.T("authentication-error")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	stt := status.New(codes.Unauthenticated, detail)

	preconditionFailure := &errdetails.PreconditionFailure{}
	for _, detailErr := range a.detailErrs {
		authenticationErrTrans, err := detailErr.Translate(trans)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		violation := &errdetails.PreconditionFailure_Violation{
			Type:        "AUTHENTICATION",
			Description: authenticationErrTrans,
		}

		preconditionFailure.Violations = append(preconditionFailure.Violations, violation)
	}

	stt, err = stt.WithDetails(preconditionFailure)
	return stt, errors.WithStack(err)
}

func (a authenticationErrors) Message(trans ut.Translator) (string, error) {
	builder := new(strings.Builder)

	detail, err := trans.T("authentication-error")
	if err != nil {
		return "", errors.WithStack(err)
	}
	builder.WriteString(detail)
	for _, detailErr := range a.detailErrs {
		builder.WriteString("\n")
		authenticationErrTrans, err := detailErr.Translate(trans)
		if err != nil {
			return "", errors.WithStack(err)
		}
		builder.WriteString(authenticationErrTrans)
	}

	return builder.String(), nil
}
