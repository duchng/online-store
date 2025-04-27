package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"

	"store-management/pkg/atomicity"
	"store-management/pkg/configuration"
	"store-management/pkg/environment"
	"store-management/pkg/shutdown"
)

type (
	GetDbFunc func(ctx context.Context) bun.IDB
)

type NewDbOptionContainer struct {
	PostDatabaseInitHook []func(*bun.DB) error
}

type NewDbOpt func(*NewDbOptionContainer)

// WithPostDatabaseInitHook run a hook after the database is initialized.
func WithPostDatabaseInitHook(hook func(*bun.DB) error) NewDbOpt {
	return func(container *NewDbOptionContainer) {
		container.PostDatabaseInitHook = append(container.PostDatabaseInitHook, hook)
	}
}

func New(
	cfg configuration.DbConfig,
	env environment.Environment,
	tasks *shutdown.Tasks,
	migrationSource fs.FS,
	opts ...NewDbOpt,
) (GetDbFunc, *atomicity.DbAtomicExecutor, error) {
	emptyAtomicExecutor := &atomicity.DbAtomicExecutor{}
	completeDsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?binary_parameters=yes&sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port,
		cfg.DbName,
	)
	conn, err := sql.Open("postgres", completeDsn)
	if err != nil {
		return nil, emptyAtomicExecutor, err
	}
	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(25)
	conn.SetConnMaxIdleTime(5 * time.Minute)
	conn.SetConnMaxLifetime(2 * time.Hour)

	db := bun.NewDB(conn, pgdialect.New())
	if env.IsLocal() {
		db.AddQueryHook(
			bundebug.NewQueryHook(
				bundebug.WithEnabled(true),
				bundebug.WithVerbose(true),
			),
		)
	}
	if err := conn.Ping(); err != nil {
		return nil, emptyAtomicExecutor, err
	}

	if cfg.AutoMigrate {
		err := MigrationUp(cfg.DbName, conn, migrationSource)
		switch {
		case errors.Is(err, migrate.ErrNoChange):
			break
		case err != nil:
			return nil, emptyAtomicExecutor, err
		}
	}
	newDbOpts := NewDbOptionContainer{}
	for _, opt := range opts {
		opt(&newDbOpts)
	}
	if len(newDbOpts.PostDatabaseInitHook) > 0 {
		for _, hook := range newDbOpts.PostDatabaseInitHook {
			if err := hook(db); err != nil {
				return nil, emptyAtomicExecutor, err
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

func MigrationUp(dbName string, db *sql.DB, migrations fs.FS) error {
	iofsDriver, err := iofs.New(migrations, "migrations")
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	migrator, err := migrate.NewWithInstance("iofs", iofsDriver, dbName, driver)
	if err != nil {
		return err
	}

	return migrator.Up()
}
