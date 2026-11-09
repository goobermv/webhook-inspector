package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	"github.com/goobermv/webhook-inspector/internal/domain"
)

type IngestInput struct {
	BinID       string
	Method      string
	Path        string
	Query       map[string]any
	Headers     map[string]string
	Body        []byte
	ContentType string
	Signature   string
}

type IngestService struct {
	bins     BinRepo
	requests RequestRepo
}

func NewIngestService(bins BinRepo, requests RequestRepo) *IngestService {
	return &IngestService{bins: bins, requests: requests}
}

func (s *IngestService) Ingest(ctx context.Context, in IngestInput) (*domain.CapturedRequest, error) {
	bin, err := s.bins.Get(ctx, in.BinID)
	if err != nil {
		return nil, err
	}

	if bin.ExpiresAt != nil && bin.ExpiresAt.Before(time.Now()) {
		return nil, domain.ErrBinExpired
	}

	if in.Signature != "" {
		if !validHMAC(bin.Secret, in.Body, in.Signature) {
			return nil, domain.ErrInvalidSignature
		}
	}

	if bin.RequestCount >= bin.MaxRequests {
		return nil, domain.ErrBinFull
	}

	req := &domain.CapturedRequest{
		BinID:       in.BinID,
		Method:      in.Method,
		Path:        in.Path,
		Query:       in.Query,
		Headers:     in.Headers,
		Body:        in.Body,
		ContentType: in.ContentType,
		SizeBytes:   len(in.Body),
	}
	return s.requests.Insert(ctx, req)
}

func validHMAC(secret string, body []byte, signature string) bool {
	sig := strings.ToLower(strings.TrimSpace(signature))
	sig = strings.TrimPrefix(sig, "sha256=")
	sig = strings.TrimPrefix(sig, "sha256:")
	sig = strings.TrimPrefix(sig, "v1=")

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	want := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(want), []byte(sig))
}
