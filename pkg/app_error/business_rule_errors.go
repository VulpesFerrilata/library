package app_error

import (
	"fmt"
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/kataras/iris/v12"
	"github.com/pkg/errors"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ AppError = new(BusinessRuleErrors)

type BusinessRuleErrors []BusinessRuleError

func (bre BusinessRuleErrors) Error() string {
	builder := new(strings.Builder)

	builder.WriteString("the request has violate one or more business rules")
	builder.WriteString("\n")

	for _, businessRuleError := range bre {
		builder.WriteString(businessRuleError.Error())
		builder.WriteString("\n")
	}

	return strings.TrimSpace(builder.String())
}

func (bre BusinessRuleErrors) Problem(trans ut.Translator) (iris.Problem, error) {
	problem := iris.NewProblem()
	problem.Type("about:blank")

	problem.Status(iris.StatusUnprocessableEntity)
	detail, err := trans.T("business-rule-error")
	if err != nil {
		return nil, fmt.Errorf("%w: %s", err, "business-rule-error")
	}
	problem.Detail(detail)

	var errs []string
	for _, businessRuleError := range bre {
		businessRuleErrorTrans, err := businessRuleError.Translate(trans)
		if err != nil {
			return nil, errors.Wrap(err, "app_error.BusinessRuleErrors.Problem")
		}
		errs = append(errs, businessRuleErrorTrans)
	}
	problem.Key("errors", errs)

	return problem, nil
}

func (bre BusinessRuleErrors) Status(trans ut.Translator) (*status.Status, error) {
	detail, err := trans.T("business-rule-error")
	if err != nil {
		return nil, fmt.Errorf("%w: %s", err, "business-rule-error")
	}
	stt := status.New(codes.FailedPrecondition, detail)

	preconditionFailure := &errdetails.PreconditionFailure{}
	for _, businessRuleError := range bre {
		description, err := businessRuleError.Translate(trans)
		if err != nil {
			return nil, errors.Wrap(err, "app_error.BusinessRuleErrors.Status")
		}

		violation := &errdetails.PreconditionFailure_Violation{
			Type:        "business rule",
			Description: description,
		}

		preconditionFailure.Violations = append(preconditionFailure.Violations, violation)
	}

	stt, err = stt.WithDetails(preconditionFailure)
	return stt, errors.Wrap(err, "app_error.BusinessRuleErrors.Status")
}

func (bre BusinessRuleErrors) Message(trans ut.Translator) (string, error) {
	builder := new(strings.Builder)

	for _, fieldError := range bre {
		fieldErrorTrans, err := fieldError.Translate(trans)
		if err != nil {
			return "", errors.Wrap(err, "app_error.BusinessRuleErrors.Message")
		}
		builder.WriteString(fieldErrorTrans)
		builder.WriteString("\n")
	}

	return builder.String(), nil
}
