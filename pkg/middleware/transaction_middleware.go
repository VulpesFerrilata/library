package middleware

import (
	"context"
	"database/sql"

	"github.com/asim/go-micro/v3/server"
	"github.com/kataras/iris/v12"
	iris_context "github.com/kataras/iris/v12/context"
	"gorm.io/gorm"
)

type transactionKey struct{}

func NewTransactionMiddleware(db *gorm.DB) *TransactionMiddleware {
	return &TransactionMiddleware{
		db: db,
	}
}

type TransactionMiddleware struct {
	db *gorm.DB
}

func (t TransactionMiddleware) ServeWithTxOptions(opts ...*sql.TxOptions) iris.Handler {
	opt := new(sql.TxOptions)
	if len(opts) > 0 {
		opt = opts[0]
	}

	return func(ctx iris.Context) {
		request := ctx.Request()
		requestCtx := request.Context()
		tx := t.db.WithContext(requestCtx).Begin(opt)

		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
				panic(r)
			}

			tx.Commit()
		}()

		requestCtx = context.WithValue(requestCtx, transactionKey{}, tx)
		request.WithContext(requestCtx)
		ctx.ResetRequest(request)

		ctx.Next()

		statusCode := ctx.GetStatusCode()
		if iris_context.StatusCodeNotSuccessful(statusCode) {
			tx.Rollback()
		}
	}
}

func (t TransactionMiddleware) HandlerWrapperWithTxOptions(opts ...*sql.TxOptions) server.HandlerWrapper {
	opt := new(sql.TxOptions)
	if len(opts) > 0 {
		opt = opts[0]
	}

	return func(f server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, request server.Request, response interface{}) error {
			tx := t.db.WithContext(ctx).Begin(opt)

			defer func() {
				if r := recover(); r != nil {
					tx.Rollback()
					panic(r)
				}

				tx.Commit()
			}()

			ctx = context.WithValue(ctx, transactionKey{}, tx)

			err := f(ctx, request, response)
			if err != nil {
				tx.Rollback()
				return err
			}

			return nil
		}
	}
}

func (t TransactionMiddleware) Get(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(transactionKey{}).(*gorm.DB)
	if !ok {
		return t.db.WithContext(ctx)
	}
	return tx.WithContext(ctx)
}
