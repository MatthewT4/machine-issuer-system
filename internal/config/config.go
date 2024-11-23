package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	DbURL          string `env:"DB_URL,required"`
	ApiServerPort  uint16 `env:"API_SERVER_PORT,required" envDefault:"80"`
	AuthSecretKey  string `env:"AUTH_SECRET_KEY,required"`
	AuthTTL        int64  `env:"AUTH_TTL" envDefault:"24"` // hours
	AuthCookieName string `env:"AUTH_COOKIE_NAME" envDefault:"session_cookie"`
}

func GetConfig() (Config, error) {
	return env.ParseAs[Config]()
}
