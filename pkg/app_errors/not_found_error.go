package app_errors

import (
	"fmt"

	ut "github.com/go-playground/universal-translator"
	"github.com/kataras/iris/v12"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewNotFoundError(name string) AppError {
	return &NotFoundError{
		name: name,
	}
}

type NotFoundError struct {
	name string
}

func (nfe NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", nfe.name)
}

func (nfe NotFoundError) Problem(trans ut.Translator) (iris.Problem, error) {
	problem := iris.NewProblem()
	problem.Type("about:blank")
	problem.Status(iris.StatusNotFound)
	title, err := trans.T("not-found-error")
	if err != nil {
		return nil, err
	}
	problem.Title(title)
	detail, err := trans.T("not-found-error-detail", nfe.name)
	if err != nil {
		return nil, err
	}
	problem.Detail(detail)
	return problem, nil
}

func (nfe NotFoundError) Status(trans ut.Translator) (*status.Status, error) {
	detail, err := trans.T("not-found-error-detail", nfe.name)
	if err != nil {
		return nil, err
	}
	stt := status.New(codes.NotFound, detail)
	return stt, nil
}

func (nfe NotFoundError) Message(trans ut.Translator) (string, error) {
	return trans.T("not-found-error-detail", nfe.name)
}
