package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/thawatchai/full-stack-authentication/back-end/internal/config"
	"github.com/thawatchai/full-stack-authentication/back-end/internal/database"
	"github.com/thawatchai/full-stack-authentication/back-end/internal/repository"
	"github.com/thawatchai/full-stack-authentication/back-end/internal/worker"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	poolCfg := database.PoolConfigFrom(cfg)
	db, err := database.Connect(cfg.DatabaseURL, poolCfg)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer database.Close(db)

	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	cleanup := worker.NewTokenCleanup(refreshTokenRepo, time.Hour)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go cleanup.Run(ctx)

	slog.Info("token cleanup worker started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	cancel()
	slog.Info("token cleanup worker stopped")
}
