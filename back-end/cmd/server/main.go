package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"github.com/thawatchai/full-stack-authentication/back-end/internal/bcryptpool"
	"github.com/thawatchai/full-stack-authentication/back-end/internal/cache"
	"github.com/thawatchai/full-stack-authentication/back-end/internal/config"
	"github.com/thawatchai/full-stack-authentication/back-end/internal/database"
	"github.com/thawatchai/full-stack-authentication/back-end/internal/handlers"
	"github.com/thawatchai/full-stack-authentication/back-end/internal/middleware"
	"github.com/thawatchai/full-stack-authentication/back-end/internal/repository"
	"github.com/thawatchai/full-stack-authentication/back-end/internal/service"
	"github.com/thawatchai/full-stack-authentication/back-end/internal/telemetry"
	"github.com/thawatchai/full-stack-authentication/back-end/internal/worker"
)

const serviceName = "auth-api"

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	if cfg.IsProduction() {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	}

	ctx := context.Background()
	shutdownTracing, err := telemetry.InitTracing(ctx, serviceName, cfg.OTLPEndpoint)
	if err != nil {
		log.Fatalf("init tracing: %v", err)
	}

	poolCfg := database.PoolConfigFrom(cfg)
	db, err := database.Connect(cfg.DatabaseURL, poolCfg)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}

	readDB, err := database.Connect(cfg.DatabaseReadURL, poolCfg)
	if err != nil {
		log.Fatalf("connect read database: %v", err)
	}

	var redisClient *redis.Client
	if cfg.RedisURL != "" {
		redisClient, err = cache.NewRedisClient(cfg.RedisURL)
		if err != nil {
			log.Fatalf("connect redis: %v", err)
		}
	}

	userRepo := repository.NewUserRepository(db, readDB)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)

	var refreshStore cache.RefreshTokenStore = cache.NewPGRefreshTokenStore(refreshTokenRepo)
	var userCache cache.UserCache
	if redisClient != nil {
		refreshStore = cache.NewRedisRefreshTokenStore(redisClient, refreshTokenRepo)
		userCache = cache.NewRedisUserCache(redisClient, 10*time.Minute)
	}

	bcryptPool := bcryptpool.New(cfg.BcryptMaxConcurrent)
	authService := service.NewAuthService(cfg, userRepo, refreshTokenRepo, refreshStore, userCache, bcryptPool)
	authHandler := handlers.NewAuthHandler(authService)
	healthHandler := handlers.NewHealthHandler(db, redisClient)

	cleanupCtx, cleanupCancel := context.WithCancel(context.Background())
	defer cleanupCancel()
	go worker.NewTokenCleanup(refreshTokenRepo, time.Hour).Run(cleanupCtx)

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	app.Use(middleware.Metrics())
	app.Use(telemetry.Tracing(serviceName))

	if cfg.IsProduction() {
		app.Use(middleware.RequestLogger())
	} else {
		app.Use(fiberlogger.New())
	}
	
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSOrigin,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Request-ID",
	}))

	app.Get("/health", healthHandler.Check)
	app.Get("/ready", healthHandler.Ready)
	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	auth := app.Group("/auth")

	var authLimiter fiber.Handler = middleware.InMemoryAuthRateLimiter()
	if redisClient != nil {
		authLimiter = middleware.AuthRateLimiter(redisClient)
	}

	auth.Post("/register", authLimiter, authHandler.Register)
	auth.Post("/login", authLimiter, authHandler.Login)
	auth.Post("/refresh", authLimiter, authHandler.Refresh)

	protected := auth.Group("", middleware.Auth(authService))
	protected.Post("/logout", authHandler.Logout)
	protected.Get("/me", authHandler.Me)

	addr := fmt.Sprintf(":%s", cfg.Port)
	go func() {
		slog.Info("server starting", "addr", addr, "env", cfg.AppEnv)
		if err := app.Listen(addr); err != nil {
			slog.Error("server stopped", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server")
	cleanupCancel()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		slog.Error("server shutdown failed", "error", err)
	}

	if err := database.Close(readDB); err != nil {
		slog.Error("close read db failed", "error", err)
	}
	if err := database.Close(db); err != nil {
		slog.Error("close db failed", "error", err)
	}
	if redisClient != nil {
		_ = redisClient.Close()
	}
	if err := shutdownTracing(shutdownCtx); err != nil {
		slog.Error("tracing shutdown failed", "error", err)
	}
}
