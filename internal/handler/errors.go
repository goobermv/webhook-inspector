package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/goobermv/webhook-inspector/internal/domain"
	"github.com/goobermv/webhook-inspector/pkg/httputil"
)

func writeError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, domain.ErrBinNotFound):
		httputil.Error(w, http.StatusNotFound, "bin not found")
	case errors.Is(err, domain.ErrRequestNotFound):
		httputil.Error(w, http.StatusNotFound, "request not found")
	case errors.Is(err, domain.ErrBinExpired):
		httputil.Error(w, http.StatusGone, "bin has expired")
	case errors.Is(err, domain.ErrInvalidSignature):
		httputil.Error(w, http.StatusUnauthorized, "invalid hmac signature")
	case errors.Is(err, domain.ErrBinFull):
		httputil.Error(w, http.StatusTooManyRequests, "bin request limit reached")
	case errors.Is(err, domain.ErrBodyTooLarge):
		httputil.Error(w, http.StatusRequestEntityTooLarge, "request body exceeds limit")
	case errors.Is(err, domain.ErrInvalidTargetURL),
		errors.Is(err, domain.ErrForbiddenTarget):
		httputil.Error(w, http.StatusBadRequest, err.Error())
	default:
		slog.ErrorContext(r.Context(), "unhandled error",
			"error", err, "path", r.URL.Path, "method", r.Method)
		httputil.Error(w, http.StatusInternalServerError, "internal server error")
	}
}
