package app_error

import (
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/kataras/iris/v12"
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
	title, err := trans.T("unprocessable-entity-error")
	if err != nil {
		return nil, err
	}
	problem.Title(title)

	detail, err := trans.T("business-rule-error")
	if err != nil {
		return nil, err
	}
	problem.Detail(detail)

	var errors []string
	for _, businessRuleError := range bre {
		businessRuleErrorTrans, err := businessRuleError.Translate(trans)
		if err != nil {
			return nil, err
		}
		errors = append(errors, businessRuleErrorTrans)
	}
	problem.Key("errors", errors)

	return problem, nil
}

func (bre BusinessRuleErrors) Status(trans ut.Translator) (*status.Status, error) {
	detail, err := trans.T("business-rule-error")
	if err != nil {
		return nil, err
	}

	stt := status.New(codes.FailedPrecondition, detail)

	preconditionFailure := &errdetails.PreconditionFailure{}

	violationType, err := trans.T("business-rule")
	if err != nil {
		return nil, err
	}
	for _, businessRuleError := range bre {
		description, err := businessRuleError.Translate(trans)
		if err != nil {
			return nil, err
		}

		violation := &errdetails.PreconditionFailure_Violation{
			Type:        violationType,
			Description: description,
		}

		preconditionFailure.Violations = append(preconditionFailure.Violations, violation)
	}

	return stt.WithDetails(preconditionFailure)
}

func (bre BusinessRuleErrors) Message(trans ut.Translator) (string, error) {
	builder := new(strings.Builder)

	for _, fieldError := range bre {
		fieldErrorTrans, err := fieldError.Translate(trans)
		if err != nil {
			return "", err
		}
		builder.WriteString(fieldErrorTrans)
		builder.WriteString("\n")
	}

	return builder.String(), nil
}
