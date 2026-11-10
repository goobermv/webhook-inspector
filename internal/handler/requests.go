package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/goobermv/webhook-inspector/internal/service"
	"github.com/goobermv/webhook-inspector/pkg/httputil"
)

const (
	defaultLimit = 50
	maxLimit     = 200
)

type RequestHandler struct {
	svc *service.RequestService
}

func NewRequestHandler(svc *service.RequestService) *RequestHandler {
	return &RequestHandler{svc: svc}
}

func (h *RequestHandler) ListByBin(w http.ResponseWriter, r *http.Request) {
	binID, ok := parseUUID(w, chi.URLParam(r, "id"))
	if !ok {
		return
	}

	limit := parseIntQuery(r, "limit", defaultLimit, 1, maxLimit)
	offset := parseIntQuery(r, "offset", 0, 0, 1_000_000)

	reqs, err := h.svc.ListByBin(r.Context(), binID, limit, offset)
	if err != nil {
		writeError(w, r, err)
		return
	}
	httputil.JSON(w, http.StatusOK, reqs)
}

func (h *RequestHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUID(w, chi.URLParam(r, "id"))
	if !ok {
		return
	}

	req, err := h.svc.Get(r.Context(), id)
	if err != nil {
		writeError(w, r, err)
		return
	}
	httputil.JSON(w, http.StatusOK, req)
}

func (h *RequestHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUID(w, chi.URLParam(r, "id"))
	if !ok {
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		writeError(w, r, err)
		return
	}
	httputil.NoContent(w)
}

func parseIntQuery(r *http.Request, key string, def, min, max int) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < min || n > max {
		return def
	}
	return n
}
