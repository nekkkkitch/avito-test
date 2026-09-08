package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type ServerConfig struct {
	Address         string        `env:"ADDRESS" env-default:":8081"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" env-default:"10s"`
}

type MainServerConfig struct {
	URL string `env:"URL" env-default:"http://localhost:8080"`
}

type Config struct {
	Server     ServerConfig     `env-prefix:"MARKET_"`
	MainServer MainServerConfig `env-prefix:"MAIN_SERVER_"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("read env: %w", err)
	}
	return &cfg, nil
}
