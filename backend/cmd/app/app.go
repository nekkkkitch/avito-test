package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"backend/config"
	repo "backend/rdb"
	rdb "backend/rdb/postgres"
	"backend/server"
	"backend/services/cart"
	"backend/services/external"
	"backend/services/products"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	cfg    *config.Config
	fiber  *fiber.App
	server *server.Server
	pool   *pgxpool.Pool
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	m, err := migrate.New(
		"file://./migrations",
		cfg.Postgres.DSN(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare migrations: %w", err)
	}

	err = m.Up()

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	slog.Info("migrations applied successfully")

	pool, err := pgxpool.New(ctx, cfg.Postgres.DSN())
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		pool.Close()

		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	queries := rdb.New(pool)
	cartRepo := repo.NewCartRepo(queries)
	externalRepo := repo.NewExternalRepo(queries)
	productsRepo := repo.NewProductsRepo(queries)

	cartSvc := cart.New(&cartRepo)
	externalSvc := external.New(&externalRepo)
	productsSvc := products.New(&productsRepo)

	app := fiber.New()
	srv := server.New(app, productsSvc, cartSvc, externalSvc)
	srv.SetupRoutes()

	return &App{
		cfg:    cfg,
		fiber:  app,
		server: srv,
		pool:   pool,
	}, nil
}

// Run блокируется на прослушивании HTTP до остановки сервера
func (a *App) Run() error {
	return a.fiber.Listen(a.cfg.Server.Address)
}

// Shutdown гасит HTTP-сервер с дедлайном ctx и закрывает внешние ресурсы
func (a *App) Shutdown(ctx context.Context) error {
	err := a.fiber.ShutdownWithContext(ctx)

	a.pool.Close()

	if err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}

	return nil
}
