package atomicity

import (
	"context"

	"github.com/uptrace/bun"
)

type ContextKey string

const TxKey ContextKey = "transactionInstance"

type DbAtomicExecutor struct {
	DB *bun.DB
}

// Execute executes the given function in a transaction.
// if the parent context already contains a transaction, it will be used.
func (e *DbAtomicExecutor) Execute(parentCtx context.Context, executeFunc func(ctx context.Context) error) (err error) {
	if tx := ContextGetTx(parentCtx); tx.Tx != nil {
		return executeFunc(parentCtx)
	}
	return e.DB.RunInTx(
		parentCtx, nil, func(ctx context.Context, tx bun.Tx) error {
			return executeFunc(ContextSetTx(ctx, tx))
		},
	)
}

func ContextSetTx(ctx context.Context, tx bun.Tx) context.Context {
	return context.WithValue(ctx, TxKey, tx)
}

func ContextGetTx(ctx context.Context) bun.Tx {
	if tx, ok := ctx.Value(TxKey).(bun.Tx); ok {
		return tx
	}
	return bun.Tx{}
}

func ContextWithoutTx(ctx context.Context) context.Context {
	return context.WithValue(ctx, TxKey, nil)
}
