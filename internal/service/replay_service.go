package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
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
	requests     RequestRepo
	replays      ReplayRepo
	allowPrivate bool
}

func NewReplayService(requests RequestRepo, replays ReplayRepo, allowPrivate bool) *ReplayService {
	return &ReplayService{
		requests:     requests,
		replays:      replays,
		allowPrivate: allowPrivate,
	}
}

func (s *ReplayService) Replay(ctx context.Context, requestID, targetURL string) (*domain.Replay, error) {
	if err := validateTarget(targetURL, s.allowPrivate); err != nil {
		return nil, err
	}

	orig, err := s.requests.GetByID(ctx, requestID)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, orig.Method, targetURL, bytes.NewReader(orig.Body))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidTargetURL, err)
	}
	for k, v := range orig.Headers {
		httpReq.Header.Set(k, v)
	}

	replay := &domain.Replay{
		RequestID: orig.ID,
		TargetURL: targetURL,
	}

	client := &http.Client{Timeout: replayTimeout}

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

func validateTarget(rawURL string, allowPrivate bool) error {
	u, err := url.Parse(rawURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return domain.ErrInvalidTargetURL
	}

	if allowPrivate {
		slog.Warn("replay to potentially private target allowed",
			"url", rawURL, "reason", "ALLOW_PRIVATE_TARGETS=true")
		return nil
	}

	host := u.Hostname()
	ips, err := net.LookupIP(host)
	if err != nil {
		return domain.ErrInvalidTargetURL
	}
	for _, ip := range ips {
		if isBlockedIP(ip) {
			return domain.ErrForbiddenTarget
		}
	}
	return nil
}

func isBlockedIP(ip net.IP) bool {
	return ip.IsLoopback() ||
		ip.IsPrivate() ||
		ip.IsLinkLocalUnicast() ||
		ip.IsLinkLocalMulticast() ||
		ip.IsUnspecified() ||
		ip.IsMulticast() ||
		ip.Equal(net.IPv4bcast)
}
