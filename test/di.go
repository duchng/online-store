package integration

import (
	"context"
	"log/slog"
	"os"

	"github.com/samber/do/v2"
	"github.com/uptrace/bun"

	"store-management/internal/assets"
	"store-management/pkg/atomicity"
	"store-management/pkg/configuration"
	"store-management/pkg/database"
	"store-management/pkg/environment"
	"store-management/pkg/shutdown"
)

type TestInjectorOpt func(i do.Injector)

func WithDb(db *bun.DB) TestInjectorOpt {
	return func(i do.Injector) {
		do.ProvideValue[*bun.DB](i, db)
	}
}

func NewInjector[T any](i do.Injector, opts ...TestInjectorOpt) do.Injector {
	cfg, _ := configuration.InitConfig[T](assets.EmbeddedFiles)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	do.ProvideValue(i, logger)
	do.ProvideValue[T](i, cfg)
	do.ProvideValue(i, environment.Development)
	do.ProvideValue[*shutdown.Tasks](i, &shutdown.Tasks{})
	do.Provide[atomicity.AtomicExecutor](
		i, func(injector do.Injector) (atomicity.AtomicExecutor, error) {
			db := do.MustInvoke[*bun.DB](i)
			return &atomicity.DbAtomicExecutor{DB: db}, nil
		},
	)
	do.Provide[database.GetDbFunc](
		i, func(injector do.Injector) (database.GetDbFunc, error) {
			db := do.MustInvoke[*bun.DB](i)
			return func(ctx context.Context) bun.IDB {
				if tx := atomicity.ContextGetTx(ctx); tx.Tx != nil {
					return tx
				}
				return db
			}, nil
		},
	)
	for _, opt := range opts {
		opt(i)
	}
	return i
}
