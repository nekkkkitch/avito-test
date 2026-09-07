package config

import (
	"fmt"
	"net/url"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type ServerConfig struct {
	Address         string        `env:"ADDRESS" env-default:":8080"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" env-default:"10s"`
}

type PostgresConfig struct {
	Host     string `env:"HOST" env-default:"localhost"`
	User     string `env:"USER" env-default:"postgres"`
	Password string `env:"PASSWORD"`
	Database string `env:"DB" env-default:"postgres"`
	SSLMode  string `env:"SSLMODE" env-default:"disable"`
	Port     int    `env:"PORT" env-default:"5432"`
}

type Config struct {
	Postgres PostgresConfig `env-prefix:"POSTGRES_"`
	Server   ServerConfig   `env-prefix:"SERVER_"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("read env: %w", err)
	}

	return &cfg, nil
}

func (p *PostgresConfig) DSN() string {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(p.User, p.Password),
		Host:   fmt.Sprintf("%s:%d", p.Host, p.Port),
		Path:   "/" + p.Database,
	}

	q := u.Query()
	q.Set("sslmode", p.SSLMode)
	u.RawQuery = q.Encode()

	return u.String()
}
