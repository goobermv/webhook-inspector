package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunRetention(ctx context.Context, db *pgxpool.Pool, olderThan time.Duration) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_, err := db.Exec(ctx,
				`DELETE FROM requests WHERE received_at < now() - $1::interval`,
				olderThan.String(),
			)
			if err != nil {
				slog.Error("retention failed", "err", err)
			}
		}
	}
}
