package app

import (
	"context"
	"errors"
	"fmt"
	"os"

	"market/config"
	repo "market/rdb"
	"market/server"
	"market/services/external"
	"market/services/orders"
	"market/services/products"

	"github.com/gofiber/fiber/v3"
)

type App struct {
	cfg    *config.Config
	fiber  *fiber.App
	server *server.Server
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	_ = ctx

	storeRepo := repo.NewMemoryRepo()
	externalSvc := external.New(storeRepo)
	productsSvc := products.New(storeRepo)
	ordersSvc := orders.New(storeRepo)

	fiberApp := fiber.New()
	srv := server.New(fiberApp, productsSvc, externalSvc, ordersSvc, cfg.MainServer.URL)
	srv.SetupRoutes()

	return &App{
		cfg:    cfg,
		fiber:  fiberApp,
		server: srv,
	}, nil
}

func (a *App) Run() error {
	return a.fiber.Listen(a.cfg.Server.Address)
}

func (a *App) Shutdown(ctx context.Context) error {
	err := a.fiber.ShutdownWithContext(ctx)
	if err != nil && !errors.Is(err, os.ErrClosed) {
		return fmt.Errorf("shutdown: %w", err)
	}
	return nil
}
