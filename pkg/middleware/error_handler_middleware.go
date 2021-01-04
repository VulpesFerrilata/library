package middleware

import (
	"context"

	"github.com/VulpesFerrilata/library/pkg/app_error"
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
		if grpcErr, ok := errors.Cause(err).(app_error.GrpcError); ok {
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

	if sttErr, ok := errors.Cause(err).(app_error.StatusError); ok {
		stt := sttErr.GRPCStatus()
		ehm.ErrorHandler(ctx, app_error.NewStatusError(stt))
		return
	}

	if businessRuleErr, ok := errors.Cause(err).(app_error.BusinessRuleError); ok {
		businessRuleErrs := make(app_error.BusinessRuleErrors, 0)
		businessRuleErrs = append(businessRuleErrs, businessRuleErr)
		ehm.ErrorHandler(ctx, businessRuleErrs)
		return
	}

	if webErr, ok := errors.Cause(err).(app_error.WebError); ok {
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
