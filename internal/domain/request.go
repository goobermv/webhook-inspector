package domain

import "time"

type CapturedRequest struct {
	ID          string            `json:"id"`
	BinID       string            `json:"bin_id"`
	Method      string            `json:"method"`
	Path        string            `json:"path"`
	Query       map[string]any    `json:"query"`
	Headers     map[string]string `json:"headers"`
	Body        []byte            `json:"body,omitempty"`
	ContentType string            `json:"content_type"`
	SizeBytes   int               `json:"size_bytes"`
	ReceivedAt  time.Time         `json:"received_at"`
}
