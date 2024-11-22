package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	DbURL         string `env:"DB_URL"`
	ApiServerPort uint16 `env:"API_SERVER_PORT"`
	Auth          Auth   `env:"AUTH"`
}

func GetConfig() (Config, error) {
	return env.ParseAs[Config]()
}

type Auth struct {
	SecretKey string `env:"SECRET_KEY"`
	TTL       string `env:"TTL"`
}
