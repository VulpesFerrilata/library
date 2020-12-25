package app_errors

import (
	"fmt"

	ut "github.com/go-playground/universal-translator"
	"github.com/kataras/iris/v12"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewNotFoundError(name string) AppError {
	return &notFoundError{
		name: name,
	}
}

type notFoundError struct {
	name string
}

func (nfe notFoundError) Error() string {
	return fmt.Sprintf("%s not found", nfe.name)
}

func (nfe notFoundError) Problem(trans ut.Translator) (iris.Problem, error) {
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

func (nfe notFoundError) Status(trans ut.Translator) (*status.Status, error) {
	detail, err := trans.T("not-found-error-detail", nfe.name)
	if err != nil {
		return nil, err
	}
	stt := status.New(codes.NotFound, detail)
	return stt, nil
}

func (nfe notFoundError) Message(trans ut.Translator) (string, error) {
	return trans.T("not-found-error-detail", nfe.name)
}
