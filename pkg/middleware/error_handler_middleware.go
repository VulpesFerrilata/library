package middleware

import (
	"context"
	"fmt"

	"github.com/VulpesFerrilata/library/pkg/app_error"
	"github.com/asim/go-micro/v3/server"
	"github.com/kataras/iris/v12"
	"github.com/pkg/errors"
	"google.golang.org/grpc/status"
)

type statusError interface {
	GRPCStatus() *status.Status
}

func NewErrorHandlerMiddleware(translatorMiddleware *TranslatorMiddleware) *ErrorHandlerMiddleware {
	return &ErrorHandlerMiddleware{
		translatorMiddleware: translatorMiddleware,
	}
}

type ErrorHandlerMiddleware struct {
	translatorMiddleware *TranslatorMiddleware
}

func (ehm ErrorHandlerMiddleware) Serve(ctx iris.Context) {
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
}

func (ehm ErrorHandlerMiddleware) HandlerWrapper(f server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		err := f(ctx, req, rsp)
		if grpcErr, ok := err.(app_error.GrpcError); ok {
			trans := ehm.translatorMiddleware.Get(ctx)
			stt, err := grpcErr.Status(trans)
			if err != nil {
				return err
			}
			return stt.Err()
		}
		return err
	}
}

func (ehm ErrorHandlerMiddleware) ErrorHandler(ctx iris.Context, err error) {
	if err == nil {
		return
	}

	if sttErr, ok := err.(statusError); ok {
		stt := sttErr.GRPCStatus()
		err = app_error.NewStatusError(stt)
	}

	if businessRuleErr, ok := err.(app_error.BusinessRuleError); ok {
		err = app_error.NewBusinessRuleErrors(businessRuleErr)
	}

	if webErr, ok := err.(app_error.WebError); ok {
		trans := ehm.translatorMiddleware.Get(ctx.Request().Context())
		problem, err := webErr.Problem(trans)
		if err == nil {
			err = problem
		}
	}

	problem, ok := err.(iris.Problem)
	if !ok {
		problem = iris.NewProblem()
		problem.Type("about:blank")
		problem.Status(iris.StatusInternalServerError)
		problem.Title("internal server error")
		problem.Detail(fmt.Sprintf("%+v", err))
	}

	ctx.Problem(problem)
	ctx.StopExecution()
}
