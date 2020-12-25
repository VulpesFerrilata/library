package middleware

import (
	"context"

	"github.com/VulpesFerrilata/library/pkg/app_errors"
	"github.com/kataras/iris/v12"
	"github.com/micro/go-micro/v2/server"
	"github.com/pkg/errors"
)

func NewErrorHandlerMiddleware(translatorMiddleware *TranslatorMiddleware) *ErrorHandlerMiddleware {
	return &ErrorHandlerMiddleware{
		translatorMiddleware: translatorMiddleware,
	}
}

type ErrorHandlerMiddleware struct {
	translatorMiddleware *TranslatorMiddleware
}

func (ehm ErrorHandlerMiddleware) HandlerWrapper(f server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		err := f(ctx, req, rsp)
		if err, ok := errors.Cause(err).(app_errors.GrpcError); ok {
			trans := ehm.translatorMiddleware.Get(ctx)
			stt, err := err.Status(trans)
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

	if se, ok := errors.Cause(err).(app_errors.StatusError); ok {
		stt := se.GRPCStatus()
		ehm.ErrorHandler(ctx, app_errors.NewStatusError(stt))
		return
	}

	if webErr, ok := errors.Cause(err).(app_errors.WebError); ok {
		trans := ehm.translatorMiddleware.Get(ctx.Request().Context())
		problem, err := webErr.Problem(trans)
		if err != nil {
			ehm.ErrorHandler(ctx, err)
			return
		}
		ctx.Problem(problem)
		return
	}

	problem := iris.NewProblem()
	problem.Type("about:blank")
	problem.Status(iris.StatusInternalServerError)
	problem.Title("internal server error")
	problem.Detail(err.Error())
	ctx.Problem(problem)
}
