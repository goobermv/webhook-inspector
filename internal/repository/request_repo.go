package repository

import (
	"context"
	"errors"

	"github.com/goobermv/webhook-inspector/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RequestRepository struct {
	db *pgxpool.Pool
}

func NewRequestRepository(db *pgxpool.Pool) *RequestRepository {
	return &RequestRepository{
		db: db,
	}
}

func (r *RequestRepository) Insert(ctx context.Context, req *domain.CapturedRequest) (*domain.CapturedRequest, error) {
	const q = `
		INSERT INTO requests (bin_id, method, path, query, headers, body, content_type, size_bytes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, received_at`
	row := r.db.QueryRow(ctx, q,
		req.BinID, req.Method, req.Path, req.Query, req.Headers,
		req.Body, req.ContentType, req.SizeBytes,
	)
	err := row.Scan(&req.ID, &req.ReceivedAt)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (r *RequestRepository) ListByBin(ctx context.Context, binID string, limit, offset int) ([]*domain.CapturedRequest, error) {
	query := `
		SELECT id, bin_id, method, path, query, headers, body, content_type, size_bytes, received_at
		FROM requests
		WHERE bin_id = $1 
		ORDER BY received_at DESC 
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, binID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]*domain.CapturedRequest, 0)
	for rows.Next() {
		c := &domain.CapturedRequest{}
		err := rows.Scan(&c.ID, &c.BinID, &c.Method, &c.Path, &c.Query, &c.Headers,
			&c.Body, &c.ContentType, &c.SizeBytes, &c.ReceivedAt,
		)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (r *RequestRepository) GetByID(ctx context.Context, ID string) (*domain.CapturedRequest, error) {
	query := `
		SELECT id, bin_id, method, path, query, headers, body, content_type, size_bytes, received_at
		FROM requests
		WHERE id = $1
	`
	c := &domain.CapturedRequest{}
	err := r.db.QueryRow(ctx, query, ID).Scan(&c.ID, &c.BinID, &c.Method, &c.Path, &c.Query, &c.Headers,
		&c.Body, &c.ContentType, &c.SizeBytes, &c.ReceivedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrRequestNotFound
		}
		return nil, err
	}
	return c, nil
}

func (r *RequestRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM requests WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrRequestNotFound
	}
	return nil
}
