package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunRetention(ctx context.Context, db *pgxpool.Pool, olderThan time.Duration, logger *slog.Logger) {
	if olderThan <= 0 {
		logger.Warn("retention disabled", "reason", "non-positive duration")
		return
	}

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	logger.Info("retention worker started", "older_than", olderThan.String())

	runOnce(ctx, db, olderThan, logger)

	for {
		select {
		case <-ctx.Done():
			logger.Info("retention worker stopped")
			return
		case <-ticker.C:
			runOnce(ctx, db, olderThan, logger)
		}
	}
}

func runOnce(ctx context.Context, db *pgxpool.Pool, olderThan time.Duration, logger *slog.Logger) {
	const q = `DELETE FROM requests WHERE received_at < now() - $1::interval`

	opCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	tag, err := db.Exec(opCtx, q, olderThan.String())
	if err != nil {
		if ctx.Err() == nil {
			logger.Error("retention delete failed", "err", err)
		}
		return
	}

	if n := tag.RowsAffected(); n > 0 {
		logger.Info("retention deleted old requests", "count", n)
	}
}
