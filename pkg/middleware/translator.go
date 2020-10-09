package middleware

import (
	"context"

	httpext "github.com/go-playground/pkg/v5/net/http"
	"github.com/go-playground/pure/v5"
	ut "github.com/go-playground/universal-translator"
	"github.com/kataras/iris/v12"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/server"
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

func (tm TranslatorMiddleware) Serve(ctx iris.Context) {
	r := ctx.Request()

	requestCtx := r.Context()
	languages := pure.AcceptedLanguages(r)
	trans, _ := tm.utrans.FindTranslator(languages...)
	requestCtx = context.WithValue(requestCtx, translatorKey{}, trans)
	r.WithContext(requestCtx)

	ctx.ResetRequest(r)
	ctx.Next()
}

func (tm TranslatorMiddleware) CallWrapper(f client.CallFunc) client.CallFunc {
	return func(ctx context.Context, node *registry.Node, req client.Request, rsp interface{}, opts client.CallOptions) error {
		trans := tm.Get(ctx)
		ctx = metadata.Set(ctx, httpext.AcceptedLanguage, trans.Locale())
		return f(ctx, node, req, rsp, opts)
	}
}

func (tm TranslatorMiddleware) HandlerWrapper(f server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		language, _ := metadata.Get(ctx, httpext.AcceptedLanguage)
		trans, _ := tm.utrans.FindTranslator(language)
		ctx = context.WithValue(ctx, translatorKey{}, trans)
		return f(ctx, req, rsp)
	}
}

func (tm TranslatorMiddleware) Get(ctx context.Context) ut.Translator {
	trans, found := ctx.Value(translatorKey{}).(ut.Translator)
	if !found {
		return tm.utrans.GetFallback()
	}
	return trans
}
