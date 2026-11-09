package repository

import (
	"context"

	"github.com/goobermv/webhook-inspector/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReplayRepository struct{ db *pgxpool.Pool }

func NewReplayRepository(db *pgxpool.Pool) *ReplayRepository {
	return &ReplayRepository{db: db}
}

func (r *ReplayRepository) Insert(ctx context.Context, rp *domain.Replay) (*domain.Replay, error) {
	const q = `INSERT INTO replays (request_id, target_url, status_code, response_body, error)
	           VALUES ($1, $2, $3, $4, $5)
	           RETURNING id, created_at`
	err := r.db.QueryRow(ctx, q, rp.RequestID, rp.TargetURL, rp.StatusCode, rp.ResponseBody, rp.Error).
		Scan(&rp.ID, &rp.CreatedAt)
	if err != nil {
		return nil, err
	}
	return rp, nil
}
