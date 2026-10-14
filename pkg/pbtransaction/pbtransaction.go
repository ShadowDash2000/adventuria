package pbtransaction

import (
	"context"

	"github.com/pocketbase/pocketbase/core"
)

type ctxKey struct{}

type trCallback func(ctx context.Context, txApp core.App) error

func RunInTransaction(ctx context.Context, pb core.App, fn trCallback) error {
	return GetCtxTransactionOrApp(ctx, pb).RunInTransaction(func(txApp core.App) error {
		ctx := SetCtxTransaction(ctx, txApp)
		return fn(ctx, txApp)
	})
}

func GetCtxTransactionOrApp(ctx context.Context, pb core.App) core.App {
	if tr, ok := ctx.Value(ctxKey{}).(core.App); ok {
		return tr
	}
	return pb
}

func SetCtxTransaction(ctx context.Context, pb core.App) context.Context {
	return context.WithValue(ctx, ctxKey{}, pb)
}

func AfterCommit(ctx context.Context, pb core.App, fn func()) {
	app := GetCtxTransactionOrApp(ctx, pb)

	if tx := app.TxInfo(); tx != nil {
		tx.OnComplete(func(txErr error) error {
			if txErr == nil {
				fn()
			}
			return nil
		})
		return
	}

	fn()
}
