package dbtest

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"log/slog"
	"math/rand/v2"
	"os"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"

	"store-management/pkg/database"
)

const (
	defaultPassword = "12341234"
	defaultUser     = "admin"
	dbPrefix        = "ecommerce"
)

var postgresContainer testcontainers.Container

func StartDatabase() (cleanupFunc func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	req := testcontainers.ContainerRequest{
		Image:        "postgres",
		ExposedPorts: []string{"5432/tcp"},
		HostConfigModifier: func(config *container.HostConfig) {
			config.AutoRemove = true
		},
		Env: map[string]string{
			"POSTGRES_USER":     defaultUser,
			"POSTGRES_PASSWORD": defaultPassword,
			"POSTGRES_DB":       dbPrefix,
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}
	postgres, err := testcontainers.GenericContainer(
		ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		},
	)
	if err != nil {
		slog.Error("err starting postgres container", slog.String("err", err.Error()))
		os.Exit(1)
	}
	postgresContainer = postgres
	return func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			slog.Error("err stopping postgres container", slog.String("err", err.Error()))
			os.Exit(0)
		}
	}
}

func NewDatabase(t *testing.T, migrationSource fs.FS) *bun.DB {
	// open connection to postgres instance in order to create other databases
	baseConn := connectDb(t, dbPrefix)
	defer func() {
		if err := baseConn.Close(); err != nil {
			t.Fatal("err close connection to db")
		}
	}()
	dbName := fmt.Sprintf("%s_%d", dbPrefix, rand.Int())
	if _, err := baseConn.Exec(fmt.Sprintf("CREATE DATABASE %s;", dbName)); err != nil {
		t.Fatal("err creating postgres database")
	}

	subConn := connectDb(t, dbName)

	// apply migrations to the newly created database
	if err := database.MigrationUp(dbName, subConn, migrationSource); err != nil {
		t.Fatal("err connect postgres database", err)
	}
	// connect to new database
	db := bun.NewDB(subConn, pgdialect.New())
	db.AddQueryHook(
		bundebug.NewQueryHook(
			bundebug.WithEnabled(true),
			bundebug.WithVerbose(true),
		),
	)
	return db
}

func connectDb(t *testing.T, dbName string) *sql.DB {
	ctx := context.Background()
	hostIP, err := postgresContainer.Host(ctx)
	if err != nil {
		t.Fatal("err get host ip from container")
	}
	mappedPort, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatal("err get mapped port from container")
	}
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		hostIP, mappedPort.Port(),
		defaultUser, defaultPassword, dbName,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatal("err connect postgres database")
	}
	return db
}
