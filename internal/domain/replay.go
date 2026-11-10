package domain

import "time"

type Replay struct {
	ID           string    `json:"id"`
	RequestID    string    `json:"request_id"`
	TargetURL    string    `json:"target_url"`
	StatusCode   *int      `json:"status_code,omitempty"`
	ResponseBody []byte    `json:"response_body,omitempty"`
	Error        string    `json:"error,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}
