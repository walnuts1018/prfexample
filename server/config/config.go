package config

import (
	"log/slog"
	"reflect"

	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	LogLevel slog.Level `env:"LOG_LEVEL"`
	LogType  LogType    `env:"LOG_TYPE" envDefault:"json"`

	Server         Server         `envPrefix:"SERVER_"`
	UserHandleSalt UserHandleSalt `env:"USER_HANDLE_SALT,required"`
	ScyllaDB 	 ScyllaDB        `envPrefix:"SCYLLADB_"`
}

func Load() (Config, error) {
	var cfg Config
	if err := env.ParseWithOptions(&cfg, env.Options{
		FuncMap: map[reflect.Type]env.ParserFunc{
			reflect.TypeFor[slog.Level](): returnAny(ParseLogLevel),
			reflect.TypeFor[LogType]():    returnAny(ParseLogType),
		},
	}); err != nil {
		return Config{}, err
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func returnAny[T any](f func(v string) (t T, err error)) func(v string) (any, error) {
	return func(v string) (any, error) {
		t, err := f(v)
		return any(t), err
	}
}
