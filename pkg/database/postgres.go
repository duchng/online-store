package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"

	"store-management/pkg/atomicity"
	"store-management/pkg/configuration"
	"store-management/pkg/shutdown"
)

func NewPostgres(
	ctx context.Context,
	cfg configuration.DbConfig,
	debug bool,
	tasks *shutdown.Tasks,
	migrationSource fs.FS,
	opts ...NewDbOpt,
) (GetDbFunc, *atomicity.DbAtomicExecutor, error) {
	completeDsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?binary_parameters=yes&sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port,
		cfg.DbName,
	)
	conn, err := sql.Open("postgres", completeDsn)
	if err != nil {
		return nil, nil, err
	}
	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(25)
	conn.SetConnMaxIdleTime(5 * time.Minute)
	conn.SetConnMaxLifetime(2 * time.Hour)

	if err := conn.PingContext(ctx); err != nil {
		return nil, nil, err
	}

	db := bun.NewDB(conn, pgdialect.New())

	if debug {
		db.AddQueryHook(
			bundebug.NewQueryHook(
				bundebug.WithEnabled(true),
				bundebug.WithVerbose(true),
			),
		)
	}

	if cfg.AutoMigrate {
		err := MigrationUp(cfg.DbName, conn, migrationSource)
		switch {
		case errors.Is(err, migrate.ErrNoChange):
			break
		case err != nil:
			return nil, nil, err
		}
	}

	newDbOpts := NewDbOptionContainer{}
	for _, opt := range opts {
		opt(&newDbOpts)
	}
	if len(newDbOpts.PostDatabaseInitHook) > 0 {
		for _, hook := range newDbOpts.PostDatabaseInitHook {
			if err := hook(db); err != nil {
				return nil, nil, err
			}
		}
	}

	getDbFunc := func(ctx context.Context) bun.IDB {
		if tx := atomicity.ContextGetTx(ctx); tx.Tx != nil {
			return tx
		}
		return db
	}

	tasks.AddShutdownTask(
		func(_ context.Context) error {
			return db.Close()
		},
	)

	return getDbFunc, &atomicity.DbAtomicExecutor{DB: db}, nil
}
