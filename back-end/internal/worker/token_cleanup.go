package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/thawatchai/full-stack-authentication/back-end/internal/repository"
)

type TokenCleanup struct {
	repo     *repository.RefreshTokenRepository
	interval time.Duration
}

func NewTokenCleanup(repo *repository.RefreshTokenRepository, interval time.Duration) *TokenCleanup {
	if interval <= 0 {
		interval = time.Hour
	}
	return &TokenCleanup{repo: repo, interval: interval}
}

func (w *TokenCleanup) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	w.cleanup(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.cleanup(ctx)
		}
	}
}

func (w *TokenCleanup) cleanup(ctx context.Context) {
	_ = ctx
	deleted, err := w.repo.DeleteExpired(time.Now())
	if err != nil {
		slog.Error("token cleanup failed", "error", err)
		return
	}
	if deleted > 0 {
		slog.Info("expired refresh tokens removed", "count", deleted)
	}
}
