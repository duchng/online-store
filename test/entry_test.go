package integration

import (
	"log/slog"
	"os"
	"testing"

	"store-management/pkg/dbtest"
)

func TestMain(m *testing.M) {
	slog.Info("Start integration test for online shop..")
	cleanupDb := dbtest.StartDatabase()
	code := m.Run()
	slog.Info("End integration test for online shop..")
	cleanupDb()
	os.Exit(code)
}
