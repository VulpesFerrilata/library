package middleware

import (
	"context"

	"github.com/asim/go-micro/v3/client"
	"github.com/asim/go-micro/v3/metadata"
	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/server"
	httpext "github.com/go-playground/pkg/v5/net/http"
	"github.com/go-playground/pure/v5"
	ut "github.com/go-playground/universal-translator"
	"github.com/kataras/iris/v12"
)

type translatorKey struct{}

func NewTranslatorMiddleware(utrans *ut.UniversalTranslator) *TranslatorMiddleware {
	return &TranslatorMiddleware{
		utrans: utrans,
	}
}

type TranslatorMiddleware struct {
	utrans *ut.UniversalTranslator
}

func (t TranslatorMiddleware) Serve(ctx iris.Context) {
	request := ctx.Request()

	requestCtx := request.Context()
	languages := pure.AcceptedLanguages(request)
	trans, _ := t.utrans.FindTranslator(languages...)
	requestCtx = context.WithValue(requestCtx, translatorKey{}, trans)
	request.WithContext(requestCtx)

	ctx.ResetRequest(request)
	ctx.Next()
}

func (t TranslatorMiddleware) CallWrapper(f client.CallFunc) client.CallFunc {
	return func(ctx context.Context, node *registry.Node, request client.Request, response interface{}, opts client.CallOptions) error {
		trans := t.Get(ctx)
		ctx = metadata.Set(ctx, httpext.AcceptedLanguage, trans.Locale())
		return f(ctx, node, request, response, opts)
	}
}

func (t TranslatorMiddleware) HandlerWrapper(f server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, request server.Request, response interface{}) error {
		language, _ := metadata.Get(ctx, httpext.AcceptedLanguage)
		trans, _ := t.utrans.FindTranslator(language)
		ctx = context.WithValue(ctx, translatorKey{}, trans)
		return f(ctx, request, response)
	}
}

func (t TranslatorMiddleware) Get(ctx context.Context) ut.Translator {
	trans, found := ctx.Value(translatorKey{}).(ut.Translator)
	if !found {
		return t.utrans.GetFallback()
	}
	return trans
}
