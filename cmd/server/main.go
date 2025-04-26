package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/dyleme/template/internal/config"
	exampleHandler "github.com/dyleme/template/internal/handler/example"
	"github.com/dyleme/template/internal/httpserver"
	exampleRepository "github.com/dyleme/template/internal/repository/example"
	exampleService "github.com/dyleme/template/internal/service/example"
	"github.com/dyleme/template/pkg/log"
	"github.com/dyleme/template/pkg/log/slogpretty"
	"github.com/dyleme/template/pkg/sqldatabase"
	"github.com/dyleme/template/pkg/txmanager"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	logger := setupLogger(cfg.Env)
	ctx := log.InCtx(context.Background(), logger)
	ctx, _ = signal.NotifyContext(ctx, os.Interrupt)

	db, err := sqldatabase.NewPGX(ctx, cfg.Database.ConnectionString())
	if err != nil {
		panic(err)
	}

	txGetter := txmanager.NewGetter(db)
	txManager := txmanager.NewManager(db)

	exampleRepo := exampleRepository.NewRepository(txGetter)
	exampleSrv := exampleService.NewService(exampleRepo, txManager)
	exmpleHnldr := exampleHandler.NewHandler(exampleSrv)

	router := httpserver.Route(exmpleHnldr)
	server := httpserver.New(router, cfg.Server)

	err = server.Run(ctx)
	if err != nil {
		panic(err)
	}
}

const (
	localEnv = "local"
	devEnv   = "dev"
	prodEnv  = "prod"
)

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case localEnv:
		prettyHandler := slogpretty.NewHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}) //nolint:exhaustruct //no need to set this params
		logger = slog.New(prettyHandler)
	case devEnv:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})) //nolint:exhaustruct //no need to set this params
	case prodEnv:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})) //nolint:exhaustruct //no need to set this params
	default:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})) //nolint:exhaustruct //no need to set this params
	}

	return logger
}
