package handler

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/goobermv/webhook-inspector/internal/domain"
	"github.com/goobermv/webhook-inspector/internal/service"
	"github.com/goobermv/webhook-inspector/pkg/httputil"
)

type IngestHandler struct {
	svc         *service.IngestService
	maxBodySize int64
}

func NewIngestHandler(svc *service.IngestService, maxBodySize int64) *IngestHandler {
	return &IngestHandler{svc: svc, maxBodySize: maxBodySize}
}

func (h *IngestHandler) Ingest(w http.ResponseWriter, r *http.Request) {
	binID := chi.URLParam(r, "binID")
	if _, err := uuid.Parse(binID); err != nil {
		httputil.Error(w, http.StatusNotFound, "endpoint not found")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, h.maxBodySize)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			writeError(w, r, domain.ErrBodyTooLarge)
			return
		}
		httputil.Error(w, http.StatusBadRequest, "failed to read body")
		return
	}

	in := service.IngestInput{
		BinID:       binID,
		Method:      r.Method,
		Path:        r.URL.Path,
		Query:       flattenQuery(r.URL.Query()),
		Headers:     flattenHeaders(r.Header),
		Body:        body,
		ContentType: r.Header.Get("Content-Type"),
		Signature:   r.Header.Get("X-Signature"),
	}

	req, err := h.svc.Ingest(r.Context(), in)
	if err != nil {
		writeError(w, r, err)
		return
	}

	httputil.JSON(w, http.StatusOK, map[string]string{
		"id":     req.ID,
		"status": "received",
	})
}

func flattenQuery(q url.Values) map[string]any {
	out := make(map[string]any, len(q))
	for k, vs := range q {
		if len(vs) == 1 {
			out[k] = vs[0]
		} else {
			out[k] = vs
		}
	}
	return out
}

func flattenHeaders(h http.Header) map[string]string {
	out := make(map[string]string, len(h))
	for k, vs := range h {
		out[k] = strings.Join(vs, ", ")
	}
	return out
}
