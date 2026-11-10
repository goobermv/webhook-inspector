package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/goobermv/webhook-inspector/internal/config"
	"github.com/goobermv/webhook-inspector/internal/database"
	"github.com/goobermv/webhook-inspector/internal/handler"
	"github.com/goobermv/webhook-inspector/internal/repository"
	"github.com/goobermv/webhook-inspector/internal/service"
	"github.com/goobermv/webhook-inspector/internal/worker"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("db connect", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := database.Migrate(pool, cfg.MigrationsDir); err != nil {
		logger.Error("migrate", "err", err)
		os.Exit(1)
	}

	binRepo := repository.NewBinRepository(pool)
	reqRepo := repository.NewRequestRepository(pool)
	replayRepo := repository.NewReplayRepository(pool)

	binSvc := service.NewBinService(binRepo, cfg.DefaultMaxRequests)
	reqSvc := service.NewRequestService(reqRepo)
	ingestSvc := service.NewIngestService(binRepo, reqRepo)
	replaySvc := service.NewReplayService(reqRepo, replayRepo, cfg.AllowPrivateTargets)

	go worker.RunRetention(ctx, pool, time.Duration(cfg.RetentionHours)*time.Hour, logger)

	router := handler.NewRouter(handler.Deps{
		Bins:               binSvc,
		Requests:           reqSvc,
		Ingest:             ingestSvc,
		Replay:             replaySvc,
		AdminAPIKey:        cfg.AdminAPIKey,
		RateLimitPerMinute: cfg.RateLimitPerMinute,
		MaxBodyBytes:       1 << 20,
		Logger:             logger,
	})

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		logger.Info("server started", "addr", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server", "err", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown", "err", err)
	}
}
