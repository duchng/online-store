package main

import (
	"log/slog"
	"os"
	"runtime/debug"

	_ "go.uber.org/automaxprocs"

	"store-management/internal/app"
	"store-management/pkg/shutdown"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	tasks, _ := shutdown.NewShutdownTasks(logger)
	defer func() {
		tasks.Wait(recover())
	}()
	err := app.Run(logger, tasks)
	if err != nil {
		trace := debug.Stack()
		logger.Error("cannot start application", slog.String("error", err.Error()), slog.String("stack", string(trace)))
		os.Exit(1)
	}
}
