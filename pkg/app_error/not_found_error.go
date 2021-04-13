package app_error

import (
	"fmt"

	ut "github.com/go-playground/universal-translator"
	"github.com/kataras/iris/v12"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func IsNotFoundError(err error) bool {
	_, ok := err.(notFoundError)
	return ok
}

func NewNotFoundError(name string) AppError {
	return &notFoundError{
		name: name,
	}
}

type notFoundError struct {
	name string
}

func (n notFoundError) Error() string {
	return fmt.Sprintf("%s not found", n.name)
}

func (n notFoundError) Problem(trans ut.Translator) (iris.Problem, error) {
	problem := iris.NewProblem()
	problem.Type("about:blank")
	problem.Status(iris.StatusUnprocessableEntity)

	detail, err := trans.T("not-found-error", n.name)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	problem.Detail(detail)

	return problem, nil
}

func (n notFoundError) Status(trans ut.Translator) (*status.Status, error) {
	detail, err := trans.T("not-found-error", n.name)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	stt := status.New(codes.NotFound, detail)
	return stt, nil
}

func (n notFoundError) Message(trans ut.Translator) (string, error) {
	msg, err := trans.T("not-found-error", n.name)
	return msg, errors.WithStack(err)
}
