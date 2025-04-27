package dependencies

import (
	"context"
	"crypto/x509"

	"github.com/redis/go-redis/v9"
	"github.com/samber/do/v2"

	"store-management/internal/adapter/product/http"
	productPostgres "store-management/internal/adapter/product/postgres"
	userHttp "store-management/internal/adapter/user/http"
	userPostgres "store-management/internal/adapter/user/postgres"
	userCachedAdapter "store-management/internal/adapter/user/redis"
	"store-management/internal/config"
	"store-management/internal/core/product"
	"store-management/internal/core/user"
	"store-management/pkg/atomicity"
	"store-management/pkg/database"
	"store-management/pkg/jwttoken"
)

func NewInjector() do.Injector {
	injector := do.New()

	do.Provide(injector, NewSignParser)
	do.Provide(injector, NewRedisClient)

	do.Provide(injector, NewUserPersistencePort)
	do.Provide(injector, NewUserUseCase)
	do.Provide(injector, NewUserHandler)

	do.Provide(injector, NewProductPersistencePort)
	do.Provide(injector, NewProductUseCase)
	do.Provide(injector, NewProductHandler)

	return injector
}

func NewProductPersistencePort(i do.Injector) (product.ProductPersistencePort, error) {
	getDbFunc := do.MustInvoke[database.GetDbFunc](i)
	return productPostgres.NewProductPostgresAdapter(getDbFunc), nil
}

func NewProductUseCase(i do.Injector) (product.UseCase, error) {
	productPersistencePort := do.MustInvoke[product.ProductPersistencePort](i)
	atomicExecutor := do.MustInvoke[atomicity.AtomicExecutor](i)
	return product.NewUseCase(productPersistencePort, atomicExecutor), nil
}

func NewProductHandler(i do.Injector) (*http.ProductHandler, error) {
	productUseCase := do.MustInvoke[product.UseCase](i)
	return http.NewProductHandler(productUseCase), nil
}

func NewUserPersistencePort(i do.Injector) (user.UsePersistencePort, error) {
	getDbFunc := do.MustInvoke[database.GetDbFunc](i)
	postgresAdapter := userPostgres.NewUserPostgresAdapter(getDbFunc)
	redisClient := do.MustInvoke[*redis.Client](i)
	return userCachedAdapter.NewUserRedisAdapter(postgresAdapter, redisClient), nil
}

func NewUserUseCase(i do.Injector) (user.UseCase, error) {
	userPersistencePort := do.MustInvoke[user.UsePersistencePort](i)
	signParser := do.MustInvoke[jwttoken.SignParser](i)
	atomicExecutor := do.MustInvoke[atomicity.AtomicExecutor](i)
	return user.NewUseCase(userPersistencePort, signParser, atomicExecutor), nil
}

func NewUserHandler(i do.Injector) (*userHttp.UserHandler, error) {
	userUseCase := do.MustInvoke[user.UseCase](i)
	return userHttp.NewUserHandler(userUseCase), nil
}

func NewSignParser(i do.Injector) (jwttoken.SignParser, error) {
	cfg := do.MustInvoke[config.AppConfig](i)
	signParser, err := jwttoken.New(
		x509.PureEd25519, jwttoken.WithPublicKey(cfg.Auth.JwtPublicKey),
		jwttoken.WithPrivateKey(cfg.Auth.JwtPrivateKey),
	)
	return signParser, err
}

func NewRedisClient(i do.Injector) (*redis.Client, error) {
	cfg := do.MustInvoke[config.AppConfig](i)
	var r *redis.Client
	opts := &redis.Options{
		Addr: cfg.Redis.Host,
	}
	r = redis.NewClient(opts)
	_, err := r.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return r, nil
}
