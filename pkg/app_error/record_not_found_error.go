package app_error

import (
	"fmt"

	ut "github.com/go-playground/universal-translator"
	"github.com/kataras/iris/v12"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func IsRecordNotFoundError(err error) bool {
	_, ok := errors.Cause(err).(recordNotFoundError)
	return ok
}

func NewRecordNotFoundError(name string) AppError {
	return &recordNotFoundError{
		name: name,
	}
}

type recordNotFoundError struct {
	name string
}

func (n recordNotFoundError) Error() string {
	return fmt.Sprintf("record not found: %s", n.name)
}

func (r recordNotFoundError) Problem(trans ut.Translator) (iris.Problem, error) {
	problem := iris.NewProblem()
	problem.Type("about:blank")
	problem.Status(iris.StatusUnprocessableEntity)

	detail, err := trans.T("record-not-found-error", r.name)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	problem.Detail(detail)

	return problem, nil
}

func (r recordNotFoundError) Status(trans ut.Translator) (*status.Status, error) {
	detail, err := trans.T("record-not-found-error", r.name)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	stt := status.New(codes.NotFound, detail)
	return stt, nil
}

func (r recordNotFoundError) Message(trans ut.Translator) (string, error) {
	msg, err := trans.T("record-not-found-error", r.name)
	return msg, errors.WithStack(err)
}
