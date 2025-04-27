package app

import (
	"log/slog"

	"github.com/samber/do/v2"
	"golang.org/x/sync/errgroup"

	"store-management/internal/adapter/server"
	"store-management/internal/assets"
	"store-management/internal/config"
	"store-management/internal/dependencies"
	"store-management/pkg/atomicity"
	"store-management/pkg/configuration"
	"store-management/pkg/database"
	"store-management/pkg/shutdown"
)

func Run(logger *slog.Logger, tasks *shutdown.Tasks) error {
	cfg, err := configuration.InitConfig[config.AppConfig](assets.EmbeddedFiles)
	if err != nil {
		return err
	}

	env := cfg.Env

	getDbFunc, atomicExecutor, err := database.New(
		cfg.Database, env, tasks, assets.EmbeddedFiles,
	)
	if err != nil {
		return err
	}

	injector := dependencies.NewInjector()
	do.ProvideValue(injector, logger)
	do.ProvideValue(injector, env)
	do.ProvideValue(injector, getDbFunc)
	do.ProvideValue(injector, cfg)
	do.ProvideValue[atomicity.AtomicExecutor](injector, atomicExecutor)
	do.ProvideValue(injector, tasks)
	group := errgroup.Group{}
	group.Go(
		func() error {
			return server.ServeRest(injector)
		},
	)
	return group.Wait()
}
