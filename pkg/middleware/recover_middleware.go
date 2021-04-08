package middleware

import (
	"context"
	"fmt"

	"github.com/asim/go-micro/v3/server"
	"github.com/kataras/iris/v12"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewRecoverMiddleware(translatorMiddleware *TranslatorMiddleware) *RecoverMiddleware {
	return &RecoverMiddleware{}
}

type RecoverMiddleware struct{}

func (r RecoverMiddleware) Serve(ctx iris.Context) {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = errors.New(fmt.Sprint(r))
			}

			problem := iris.NewProblem()
			problem.Type("about:blank")
			problem.Status(iris.StatusInternalServerError)
			problem.Title("internal server error")
			problem.Detail(fmt.Sprintf("%+v", err))

			ctx.Problem(problem)
			ctx.StopExecution()
		}
	}()

	ctx.Next()
}

func (r RecoverMiddleware) HandlerWrapper(f server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, request server.Request, response interface{}) (err error) {
		defer func() {
			if r := recover(); r != nil {
				_, ok := r.(error)
				if !ok {
					err = errors.New(fmt.Sprint(r))
				} else {
					err = r.(error)
				}

				stt, ok := status.FromError(errors.Cause(err))
				if !ok {
					stt = status.New(codes.Internal, fmt.Sprintf("%+v", err))
				}

				err = stt.Err()
			}
		}()

		err = f(ctx, request, response)
		return err
	}
}
