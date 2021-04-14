package middleware

import (
	"context"

	"github.com/VulpesFerrilata/library/pkg/app_error"
	"github.com/asim/go-micro/v3/server"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/pkg/errors"
	"google.golang.org/grpc/status"
)

func NewErrorHandlerMiddleware(translatorMiddleware *TranslatorMiddleware) *ErrorHandlerMiddleware {
	return &ErrorHandlerMiddleware{
		translatorMiddleware: translatorMiddleware,
	}
}

type ErrorHandlerMiddleware struct {
	translatorMiddleware *TranslatorMiddleware
}

func (e ErrorHandlerMiddleware) HandlerWrapper(f server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, request server.Request, response interface{}) error {
		err := f(ctx, request, response)
		if grpcErr, ok := errors.Cause(err).(app_error.GrpcError); ok {
			trans := e.translatorMiddleware.Get(ctx)
			stt, err := grpcErr.Status(trans)
			if err != nil {
				panic(errors.WithStack(err))
			}
			return stt.Err()
		}
		if err != nil {
			panic(errors.WithStack(err))
		}

		return nil
	}
}

func (e ErrorHandlerMiddleware) ErrorHandler(ctx iris.Context, err error) {
	if err == nil {
		return
	}

	if stt, ok := status.FromError(errors.Cause(err)); ok {
		err = app_error.NewStatusError(stt)
	}

	if validationErrs, ok := errors.Cause(err).(validator.ValidationErrors); ok {
		err = app_error.NewValidationError(validationErrs)
	}

	if businessRuleErr, ok := errors.Cause(err).(app_error.BusinessRuleError); ok {
		err = app_error.NewBusinessRuleErrors(businessRuleErr)
	}

	if webErr, ok := errors.Cause(err).(app_error.WebError); ok {
		trans := e.translatorMiddleware.Get(ctx.Request().Context())
		problem, err := webErr.Problem(trans)
		if err != nil {
			panic(errors.WithStack(err))
		}
		ctx.Problem(problem)
		ctx.StopExecution()
		return
	}

	panic(errors.WithStack(err))
}
