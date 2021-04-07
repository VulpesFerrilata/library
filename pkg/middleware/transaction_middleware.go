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

func (tm TransactionMiddleware) ServeWithTxOptions(opts ...*sql.TxOptions) iris.Handler {
	opt := new(sql.TxOptions)
	if len(opts) > 0 {
		opt = opts[0]
	}

	return func(ctx iris.Context) {
		request := ctx.Request()
		requestCtx := request.Context()
		tx := tm.db.WithContext(requestCtx).Begin(opt)

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

func (tm TransactionMiddleware) HandlerWrapperWithTxOptions(opts ...*sql.TxOptions) server.HandlerWrapper {
	opt := new(sql.TxOptions)
	if len(opts) > 0 {
		opt = opts[0]
	}

	return func(f server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, request server.Request, response interface{}) error {
			tx := tm.db.WithContext(ctx).Begin(opt)

			defer func() {
				if request := recover(); request != nil {
					tx.Rollback()
					panic(request)
				}
				tx.Commit()
			}()

			ctx = context.WithValue(ctx, transactionKey{}, tx)

			err := f(ctx, request, response)
			if err != nil {
				tx.Rollback()
				return err
			}
			tx.Commit()

			return nil
		}
	}
}

func (tm TransactionMiddleware) Get(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(transactionKey{}).(*gorm.DB)
	if !ok {
		return tm.db.WithContext(ctx)
	}
	return tx.WithContext(ctx)
}
