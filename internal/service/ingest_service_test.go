package service

import (
	"context"
	"errors"
	"testing"

	"github.com/goobermv/webhook-inspector/internal/domain"
)

type mockBinRepo struct {
	bin *domain.Bin
	err error
}

func (m *mockBinRepo) Get(ctx context.Context, id string) (*domain.Bin, error) {
	return m.bin, m.err
}
func (m *mockBinRepo) Create(ctx context.Context, name, secret string, maxReqs int) (*domain.Bin, error) {
	return nil, nil
}
func (m *mockBinRepo) List(ctx context.Context) ([]*domain.Bin, error) { return nil, nil }
func (m *mockBinRepo) Delete(ctx context.Context, id string) error     { return nil }

type mockReqRepo struct{}

func (m *mockReqRepo) Insert(ctx context.Context, req *domain.CapturedRequest) (*domain.CapturedRequest, error) {
	return req, nil
}
func (m *mockReqRepo) ListByBin(ctx context.Context, binID string, limit, offset int) ([]*domain.CapturedRequest, error) {
	return nil, nil
}
func (m *mockReqRepo) GetByID(ctx context.Context, ID string) (*domain.CapturedRequest, error) {
	return nil, nil
}
func (m *mockReqRepo) Delete(ctx context.Context, ID string) error { return nil }

func TestIngestService_Ingest(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockBins := &mockBinRepo{
			bin: &domain.Bin{ID: "123", MaxRequests: 10, RequestCount: 0},
		}
		svc := NewIngestService(mockBins, &mockReqRepo{})

		_, err := svc.Ingest(context.Background(), IngestInput{BinID: "123"})
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("bin full", func(t *testing.T) {
		mockBins := &mockBinRepo{
			bin: &domain.Bin{ID: "123", MaxRequests: 10, RequestCount: 10},
		}
		svc := NewIngestService(mockBins, &mockReqRepo{})

		_, err := svc.Ingest(context.Background(), IngestInput{BinID: "123"})
		if !errors.Is(err, domain.ErrBinFull) {
			t.Errorf("expected ErrBinFull, got %v", err)
		}
	})

	t.Run("invalid signature", func(t *testing.T) {
		mockBins := &mockBinRepo{
			bin: &domain.Bin{ID: "123", Secret: "mysecret", MaxRequests: 10},
		}
		svc := NewIngestService(mockBins, &mockReqRepo{})

		_, err := svc.Ingest(context.Background(), IngestInput{
			BinID:     "123",
			Signature: "wrong-signature",
			Body:      []byte("hello"),
		})
		if !errors.Is(err, domain.ErrInvalidSignature) {
			t.Errorf("expected ErrInvalidSignature, got %v", err)
		}
	})
}
