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

func NewBusinessRuleErrors(businessRuleErrs ...BusinessRuleError) AppError {
	return businessRuleErrors(businessRuleErrs)
}

type businessRuleErrors []BusinessRuleError

func (b businessRuleErrors) Error() string {
	builder := new(strings.Builder)

	builder.WriteString("the request has violate one or more business rules")
	for _, businessRuleErr := range b {
		builder.WriteString("\n")
		builder.WriteString(businessRuleErr.Error())
	}

	return builder.String()
}

func (b businessRuleErrors) Problem(trans ut.Translator) (iris.Problem, error) {
	problem := iris.NewProblem()
	problem.Type("about:blank")

	problem.Status(iris.StatusUnprocessableEntity)
	detail, err := trans.T("business-rule-error")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	problem.Detail(detail)

	var errs []string
	for _, businessRuleErr := range b {
		businessRuleErrorTrans, err := businessRuleErr.Translate(trans)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		errs = append(errs, businessRuleErrorTrans)
	}
	problem.Key("errors", errs)

	return problem, nil
}

func (b businessRuleErrors) Status(trans ut.Translator) (*status.Status, error) {
	detail, err := trans.T("business-rule-error")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	stt := status.New(codes.FailedPrecondition, detail)

	preconditionFailure := &errdetails.PreconditionFailure{}
	for _, businessRuleErr := range b {
		businessRuleErrTrans, err := businessRuleErr.Translate(trans)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		violation := &errdetails.PreconditionFailure_Violation{
			Type:        "BUSINESS_RULE",
			Description: businessRuleErrTrans,
		}

		preconditionFailure.Violations = append(preconditionFailure.Violations, violation)
	}

	stt, err = stt.WithDetails(preconditionFailure)
	return stt, errors.WithStack(err)
}

func (b businessRuleErrors) Message(trans ut.Translator) (string, error) {
	builder := new(strings.Builder)

	detail, err := trans.T("business-rule-error")
	if err != nil {
		return "", errors.WithStack(err)
	}
	builder.WriteString(detail)
	for _, businessRuleErr := range b {
		builder.WriteString("\n")
		businessRuleErrTrans, err := businessRuleErr.Translate(trans)
		if err != nil {
			return "", errors.WithStack(err)
		}
		builder.WriteString(businessRuleErrTrans)
	}

	return builder.String(), nil
}
