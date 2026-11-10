package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/goobermv/webhook-inspector/internal/service"
	"github.com/goobermv/webhook-inspector/pkg/httputil"
)

type ReplayHandler struct {
	svc *service.ReplayService
}

func NewReplayHandler(svc *service.ReplayService) *ReplayHandler {
	return &ReplayHandler{svc: svc}
}

type replayRequest struct {
	TargetURL string `json:"target_url"`
}

func (h *ReplayHandler) Replay(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUID(w, chi.URLParam(r, "id"))
	if !ok {
		return
	}

	var body replayRequest
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&body); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.TargetURL == "" {
		httputil.Error(w, http.StatusBadRequest, "target_url is required")
		return
	}

	replay, err := h.svc.Replay(r.Context(), id, body.TargetURL)
	if err != nil {
		writeError(w, r, err)
		return
	}
	httputil.JSON(w, http.StatusOK, replay)
}
