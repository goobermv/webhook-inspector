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

type RequestRepo interface {
	Insert(ctx context.Context, req *domain.CapturedRequest) (*domain.CapturedRequest, error)
	ListByBin(ctx context.Context, binID string, limit, offset int) ([]*domain.CapturedRequest, error)
	GetByID(ctx context.Context, ID string) (*domain.CapturedRequest, error)
}

type ReplayRepo interface {
	Insert(ctx context.Context, rp *domain.Replay) (*domain.Replay, error)
}
