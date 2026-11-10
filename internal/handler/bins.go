package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/goobermv/webhook-inspector/internal/service"
	"github.com/goobermv/webhook-inspector/pkg/httputil"
	"github.com/google/uuid"
)

type BinHandler struct {
	svc *service.BinService
}

func NewBinHandler(svc *service.BinService) *BinHandler {
	return &BinHandler{
		svc: svc,
	}
}

type createBinRequest struct {
	Name string `json:"name"`
}

func (h *BinHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createBinRequest

	if r.ContentLength > 0 {
		dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
		if err := dec.Decode(&req); err != nil {
			httputil.Error(w, http.StatusBadRequest, "invalid request body")
			return
		}
	}

	bin, err := h.svc.Create(r.Context(), req.Name)
	if err != nil {
		writeError(w, r, err)
		return
	}
	httputil.JSON(w, http.StatusCreated, bin)
}

func (h *BinHandler) List(w http.ResponseWriter, r *http.Request) {
	bins, err := h.svc.List(r.Context())
	if err != nil {
		writeError(w, r, err)
		return
	}
	httputil.JSON(w, http.StatusOK, bins)
}

func (h *BinHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUID(w, chi.URLParam(r, "id"))
	if !ok {
		return
	}

	bin, err := h.svc.Get(r.Context(), id)
	if err != nil {
		writeError(w, r, err)
		return
	}
	httputil.JSON(w, http.StatusOK, bin)
}

func (h *BinHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

func parseUUID(w http.ResponseWriter, raw string) (string, bool) {
	if _, err := uuid.Parse(raw); err != nil {
		httputil.Error(w, http.StatusBadRequest, "id must be a valid uuid")
		return "", false
	}
	return raw, true
}
