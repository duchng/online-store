package atomicity

import (
	"context"
)

// AtomicExecutor is an interface for executing functions atomically, typically in a transaction.
type AtomicExecutor interface {
	// Execute the executeFunc atomically, ie: transaction
	Execute(parentCtx context.Context, executeFunc func(ctx context.Context) error) error
}
