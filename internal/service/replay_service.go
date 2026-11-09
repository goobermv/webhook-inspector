package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/goobermv/webhook-inspector/internal/domain"
)

const (
	replayTimeout    = 10 * time.Second
	maxResponseBytes = 1 << 20
)

type ReplayService struct {
	requests RequestRepo
	replays  ReplayRepo
}

func NewReplayService(requests RequestRepo, replays ReplayRepo) *ReplayService {
	return &ReplayService{
		requests: requests,
		replays:  replays,
	}
}

func (s *ReplayService) Replay(ctx context.Context, requestID, targetURL string) (*domain.Replay, error) {
	if err := validateTarget(targetURL); err != nil {
		return nil, err
	}

	orig, err := s.requests.GetByID(ctx, requestID)
	if err != nil {
		return nil, domain.ErrRequestNotFound
	}

	httpReq, err := http.NewRequestWithContext(ctx, orig.Method, targetURL, bytes.NewReader(orig.Body))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidTargetURL, err)
	}

	for k, v := range orig.Headers {
		httpReq.Header.Set(k, v)
	}

	client := &http.Client{Timeout: replayTimeout}

	replay := &domain.Replay{
		RequestID: orig.ID,
		TargetURL: targetURL,
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		replay.Error = err.Error()
		return s.replays.Insert(ctx, replay)
	}
	defer resp.Body.Close()

	replay.StatusCode = &resp.StatusCode
	body, _ := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes))
	replay.ResponseBody = body
	return s.replays.Insert(ctx, replay)
}

func validateTarget(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		return domain.ErrInvalidTargetURL
	}

	host := u.Hostname()
	ips, err := net.LookupIP(host)
	if err != nil {
		return domain.ErrInvalidTargetURL
	}
	for _, ip := range ips {
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsUnspecified() {
			return domain.ErrForbiddenTarget
		}
	}
	return nil
}
