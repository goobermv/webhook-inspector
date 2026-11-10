package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"github.com/goobermv/webhook-inspector/internal/domain"
)

type BinService struct {
	binRepo    BinRepo
	defaultMax int
}

func NewBinService(binRepo BinRepo, defaultMax int) *BinService {
	return &BinService{
		binRepo:    binRepo,
		defaultMax: defaultMax,
	}
}

func (s *BinService) Create(ctx context.Context, name string) (*domain.Bin, error) {
	secret, err := generateSecret()
	if err != nil {
		return nil, err
	}
	return s.binRepo.Create(ctx, name, secret, s.defaultMax)
}

func (s *BinService) Get(ctx context.Context, id string) (*domain.Bin, error) {
	return s.binRepo.Get(ctx, id)
}

func (s *BinService) List(ctx context.Context) ([]*domain.Bin, error) {
	return s.binRepo.List(ctx)
}

func (s *BinService) Delete(ctx context.Context, id string) error {
	return s.binRepo.Delete(ctx, id)
}

func generateSecret() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
