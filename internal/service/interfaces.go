package service

import (
	"context"

	"github.com/goobermv/webhook-inspector/internal/domain"
)

type BinRepo interface {
	Create(ctx context.Context, name, secret string, maxReqs int) (*domain.Bin, error)
	Get(ctx context.Context, id string) (*domain.Bin, error)
	List(ctx context.Context) ([]*domain.Bin, error)
	Delete(ctx context.Context, id string) error
}
