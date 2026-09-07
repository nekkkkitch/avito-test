package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"backend/cmd/app"
	"backend/config"
)

func main() {
	if err := run(); err != nil {
		slog.Error("app stopped with error", "err", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	// ctx живёт до первого SIGINT/SIGTERM, после чего стартует graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	application, err := app.New(ctx, cfg)
	if err != nil {
		return fmt.Errorf("init app: %w", err)
	}

	// сервер слушает в отдельной горутине, ошибку Listen ловим через канал
	srvErr := make(chan error, 1)
	go func() {
		srvErr <- application.Run()
	}()

	slog.Info("server started", "address", cfg.Server.Address)

	select {
	case err = <-srvErr:
		return fmt.Errorf("run server: %w", err)
	case <-ctx.Done():
		slog.Info("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := application.Shutdown(shutdownCtx); err != nil {
		return err
	}

	slog.Info("server stopped gracefully")

	return nil
}
