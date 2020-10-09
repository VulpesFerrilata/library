package middleware

import (
	"context"

	"github.com/VulpesFerrilata/library/pkg/errors"
	"github.com/kataras/iris/v12"
	"github.com/micro/go-micro/v2/server"
	"google.golang.org/grpc/status"
)

func NewErrorMiddleware(translatorMiddleware *TranslatorMiddleware) *ErrorMiddleware {
	return &ErrorMiddleware{
		translatorMiddleware: translatorMiddleware,
		errorHandler:         newErrorHandler(translatorMiddleware),
	}
}

type ErrorMiddleware struct {
	translatorMiddleware *TranslatorMiddleware
	errorHandler         errorHandler
}

func (em ErrorMiddleware) HandlerWrapper(f server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		err := f(ctx, req, rsp)
		if err, ok := err.(errors.Error); ok {
			trans := em.translatorMiddleware.Get(ctx)
			stt := err.ToStatus(trans)
			return stt.Err()
		}
		return err
	}
}

func (em ErrorMiddleware) ErrorHandler(ctx iris.Context, err error) {
	if err == nil {
		return
	}
	em.errorHandler.handle(ctx, err)
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
