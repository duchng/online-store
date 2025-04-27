package config

import (
	"store-management/pkg/configuration"
	"store-management/pkg/environment"
)

type AppConfig struct {
	Env environment.Environment `koanf:"env"`
	Jwt struct {
		PublicKey  string `koanf:"publicKey"`
		PrivateKey string `koanf:"privateKey"`
	} `koanf:"jwt"`
	Database       configuration.DbConfig `koanf:"db"`
	ServerHttpPort int                    `koanf:"serverHttpPort"`
	Auth           struct {
		JwtPublicKey  string `koanf:"jwtPublicKey"`
		JwtPrivateKey string `koanf:"jwtPrivateKey"`
	} `koanf:"auth"`
	Redis configuration.RedisConfig `koanf:"redis"`
}
