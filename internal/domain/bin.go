package domain

import "time"

type Bin struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Secret       string     `json:"-"`
	MaxRequests  int        `json:"max_requests"`
	RequestCount int        `json:"request_count"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}
