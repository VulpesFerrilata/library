package middleware

import (
	"context"
	"database/sql"

	"github.com/kataras/iris/v12"
	"github.com/micro/go-micro/v2/server"
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

func (tm TransactionMiddleware) ServeWithTxOptions(opts *sql.TxOptions) iris.Handler {
	return func(ctx iris.Context) {
		r := ctx.Request()
		requestCtx := r.Context()
		tx := tm.db.WithContext(requestCtx).Begin(opts)

		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
			tx.Commit()
		}()

		requestCtx = context.WithValue(requestCtx, transactionKey{}, tx)
		r.WithContext(requestCtx)
		ctx.ResetRequest(r)

		ctx.Next()
	}
}

func (tm TransactionMiddleware) HandlerWrapperWithTxOptions(opts *sql.TxOptions) server.HandlerWrapper {
	return func(f server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			tx := tm.db.WithContext(ctx).Begin(opts)

			defer func() {
				if r := recover(); r != nil {
					tx.Rollback()
				}
				tx.Commit()
			}()

			ctx = context.WithValue(ctx, transactionKey{}, tx)
			return f(ctx, req, rsp)
		}
	}
}

func (tm TransactionMiddleware) Get(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(transactionKey{}).(*gorm.DB)
	if !ok {
		return tm.db
	}
	return tx
}
