package repository

import (
	"context"
	"errors"

	"github.com/goobermv/webhook-inspector/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BinRepository struct {
	db *pgxpool.Pool
}

func NewBinRepository(db *pgxpool.Pool) *BinRepository {
	return &BinRepository{
		db: db,
	}
}

func (r *BinRepository) Create(ctx context.Context, name, secret string, maxReqs int) (*domain.Bin, error) {
	query := `
		INSERT INTO bins (name, secret, max_requests)
		VALUES ($1, $2, $3)
		RETURNING id, name, secret, max_requests, request_count, expires_at, created_at
	`

	row := r.db.QueryRow(ctx, query, name, secret, maxReqs)

	b := &domain.Bin{}
	err := row.Scan(&b.ID, &b.Name, &b.Secret, &b.MaxRequests, &b.RequestCount, &b.ExpiresAt, &b.CreatedAt)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (r *BinRepository) Get(ctx context.Context, ID string) (*domain.Bin, error) {
	query := `
		SELECT id, name, secret, max_requests, request_count, expires_at, created_at
		FROM bins
		WHERE id = $1
	`
	b := &domain.Bin{}
	err := r.db.QueryRow(ctx, query, ID).Scan(&b.ID, &b.Name, &b.Secret, &b.MaxRequests, &b.RequestCount, &b.ExpiresAt, &b.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrBinNotFound
		}
		return nil, err
	}

	return b, nil
}

func (r *BinRepository) List(ctx context.Context) ([]*domain.Bin, error) {
	query := `
		SELECT id, name, secret, max_requests, request_count, expires_at, created_at
		FROM bins ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bins := make([]*domain.Bin, 0)
	for rows.Next() {
		b := &domain.Bin{}
		err := rows.Scan(&b.ID, &b.Name, &b.Secret, &b.MaxRequests,
			&b.RequestCount, &b.ExpiresAt, &b.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		bins = append(bins, b)
	}

	return bins, rows.Err()
}

func (r *BinRepository) Delete(ctx context.Context, ID string) error {
	query := `DELETE FROM bins WHERE id = $1`
	commandTag, err := r.db.Exec(ctx, query, ID)

	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return domain.ErrBinNotFound
	}
	return nil
}
