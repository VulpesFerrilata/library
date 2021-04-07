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

func (bre businessRuleErrors) Error() string {
	builder := new(strings.Builder)

	builder.WriteString("the request has violate one or more business rules")
	for _, businessRuleError := range bre {
		builder.WriteString("\n")
		builder.WriteString(businessRuleError.Error())
	}

	return builder.String()
}

func (bre businessRuleErrors) Problem(trans ut.Translator) (iris.Problem, error) {
	problem := iris.NewProblem()
	problem.Type("about:blank")

	problem.Status(iris.StatusUnprocessableEntity)
	detail, err := trans.T("business-rule-error")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	problem.Detail(detail)

	var errs []string
	for _, businessRuleError := range bre {
		businessRuleErrorTrans, err := businessRuleError.Translate(trans)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		errs = append(errs, businessRuleErrorTrans)
	}
	problem.Key("errors", errs)

	return problem, nil
}

func (bre businessRuleErrors) Status(trans ut.Translator) (*status.Status, error) {
	detail, err := trans.T("business-rule-error")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	stt := status.New(codes.FailedPrecondition, detail)

	preconditionFailure := &errdetails.PreconditionFailure{}
	for _, businessRuleError := range bre {
		description, err := businessRuleError.Translate(trans)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		violation := &errdetails.PreconditionFailure_Violation{
			Type:        "business rule",
			Description: description,
		}

		preconditionFailure.Violations = append(preconditionFailure.Violations, violation)
	}

	stt, err = stt.WithDetails(preconditionFailure)
	return stt, errors.WithStack(err)
}

func (bre businessRuleErrors) Message(trans ut.Translator) (string, error) {
	builder := new(strings.Builder)

	detail, err := trans.T("business-rule-error")
	if err != nil {
		return "", errors.WithStack(err)
	}
	builder.WriteString(detail)
	for _, businessRuleError := range bre {
		builder.WriteString("\n")
		businessRuleErrorTrans, err := businessRuleError.Translate(trans)
		if err != nil {
			return "", errors.WithStack(err)
		}
		builder.WriteString(businessRuleErrorTrans)
	}

	return builder.String(), nil
}
