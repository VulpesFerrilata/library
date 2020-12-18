package middleware

import (
	"context"

	"github.com/VulpesFerrilata/library/pkg/app_errors"
	"github.com/VulpesFerrilata/library/pkg/errors"
	"github.com/kataras/iris/v12"
	"github.com/micro/go-micro/v2/server"
	"google.golang.org/grpc/status"
)

func NewErrorMiddleware(translatorMiddleware *TranslatorMiddleware) *ErrorMiddleware {
	return &ErrorMiddleware{
		translatorMiddleware: translatorMiddleware,
	}
}

type ErrorMiddleware struct {
	translatorMiddleware *TranslatorMiddleware
}

func (em ErrorMiddleware) HandlerWrapper(f server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		err := f(ctx, req, rsp)
		if err, ok := err.(app_errors.GrpcError); ok {
			trans := em.translatorMiddleware.Get(ctx)
			stt, err := err.Status(trans)
			if err != nil {
				return err
			}
			return stt.Err()
		}
		return err
	}
}

func (em ErrorMiddleware) ErrorHandler(ctx iris.Context, err error) {
	if err == nil {
		return
	}

	if stt, ok := status.FromError(err); ok {
		err = app_errors.NewStatusError(stt)
	}

	if err, ok := err.(app_errors.WebError); ok {
		trans := em.translatorMiddleware.Get(ctx.Request().Context())
		problem, err := err.Problem(trans)
		if err != nil {
			em.ErrorHandler(ctx, err)
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

func newErrorHandler(translatorMiddleware *TranslatorMiddleware) errorHandler {
	defaultErrorHandler := newDefaultErrorHandler(translatorMiddleware)

	serverErrorHandler := newServerErrorHandler(translatorMiddleware)
	serverErrorHandler.setNext(defaultErrorHandler)

	statusErrorHandler := newStatusErrorHandler(translatorMiddleware)
	statusErrorHandler.setNext(serverErrorHandler)
	return statusErrorHandler
}

type errorHandler interface {
	setNext(h errorHandler)
	handle(ctx iris.Context, err error)
}

type baseErrorHandler struct {
	translatorMiddleware *TranslatorMiddleware
	next                 errorHandler
}

func (beh *baseErrorHandler) setNext(h errorHandler) {
	beh.next = h
}

func (beh baseErrorHandler) handle(ctx iris.Context, err error) {
	if beh.next != nil {
		beh.next.handle(ctx, err)
	}
}

func newStatusErrorHandler(translatorMiddleware *TranslatorMiddleware) errorHandler {
	return &statusErrorHandler{
		&baseErrorHandler{
			translatorMiddleware: translatorMiddleware,
		},
	}
}

type statusErrorHandler struct {
	*baseErrorHandler
}

func (seh statusErrorHandler) handle(ctx iris.Context, err error) {
	if stt, ok := status.FromError(err); ok {
		if err, ok := errors.NewStatusError(stt); ok {
			seh.baseErrorHandler.handle(ctx, err)
			return
		}
	}
	seh.baseErrorHandler.handle(ctx, err)
}

func newServerErrorHandler(translatorMiddleware *TranslatorMiddleware) errorHandler {
	return &serverErrorHandler{
		&baseErrorHandler{
			translatorMiddleware: translatorMiddleware,
		},
	}
}

type serverErrorHandler struct {
	*baseErrorHandler
}

func (seh serverErrorHandler) handle(ctx iris.Context, err error) {
	if err, ok := err.(errors.Error); ok {
		trans := seh.translatorMiddleware.Get(ctx.Request().Context())
		problem := err.ToProblem(trans)
		seh.baseErrorHandler.handle(ctx, problem)
		return
	}
	seh.baseErrorHandler.handle(ctx, err)
}

func newDefaultErrorHandler(translatorMiddleware *TranslatorMiddleware) errorHandler {
	return &defaultErrorHandler{
		&baseErrorHandler{
			translatorMiddleware: translatorMiddleware,
		},
	}
}

type defaultErrorHandler struct {
	*baseErrorHandler
}

func (deh defaultErrorHandler) handle(ctx iris.Context, err error) {
	var problem iris.Problem
	if p, ok := err.(iris.Problem); ok {
		problem = p
	} else {
		trans := deh.translatorMiddleware.Get(ctx.Request().Context())
		problem = iris.NewProblem()
		problem.Status(iris.StatusInternalServerError)
		problem.Type("about:blank")
		title, _ := trans.T("internal-error")
		problem.Title(title)
		problem.Detail(err.Error())
	}

	ctx.Problem(problem)
}
