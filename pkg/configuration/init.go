package configuration

import (
	"errors"
	iofs "io/fs"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/fs"
	"github.com/knadh/koanf/v2"

	string_helper "store-management/pkg/string-helper"
)

const (
	EnvPrefix = "APP__"
)

func InitConfig[T any](configFile iofs.FS) (T, error) {
	var config T
	k := koanf.New(".")
	configProvider := fs.Provider(configFile, "config.yaml")
	if err := k.Load(configProvider, yaml.Parser()); err != nil {
		return config, errors.New("cannot read config from file")
	}
	if err := k.Load(
		env.ProviderWithValue(
			EnvPrefix, ".", func(key string, value string) (string, any) {
				newKey := string_helper.SnakeToCamel(
					strings.Replace(
						strings.ToLower(
							strings.TrimPrefix(key, EnvPrefix),
						), "__", ".", -1,
					),
				)

				// Check if the value contains a pipe character
				// don't use comma to pretend conflict with crontab expression
				if strings.Contains(value, "|") {
					return newKey, strings.Split(value, "|")
				}

				return newKey, value
			},
		), nil,
	); err != nil {
		return config, err
	}

	if err := k.Unmarshal("", &config); err != nil {
		return config, err
	}
	return config, nil
}
