package service

import (
	"context"

	"github.com/goobermv/webhook-inspector/internal/domain"
)

type RequestService struct {
	requests RequestRepo
}

func NewRequestService(requests RequestRepo) *RequestService {
	return &RequestService{requests: requests}
}

func (s *RequestService) ListByBin(ctx context.Context, binID string, limit, offset int) ([]*domain.CapturedRequest, error) {
	return s.requests.ListByBin(ctx, binID, limit, offset)
}

func (s *RequestService) Get(ctx context.Context, id string) (*domain.CapturedRequest, error) {
	return s.requests.GetByID(ctx, id)
}

func (s *RequestService) Delete(ctx context.Context, id string) error {
	return s.requests.Delete(ctx, id)
}
