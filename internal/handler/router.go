package handler

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/goobermv/webhook-inspector/internal/middleware"
	"github.com/goobermv/webhook-inspector/internal/service"
)

type Deps struct {
	Bins     *service.BinService
	Requests *service.RequestService
	Ingest   *service.IngestService
	Replay   *service.ReplayService

	AdminAPIKey        string
	RateLimitPerMinute int
	MaxBodyBytes       int64
	Logger             *slog.Logger
}

func NewRouter(d Deps) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(middleware.Logging(d.Logger))
	r.Use(chimw.Recoverer)

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	ingestH := NewIngestHandler(d.Ingest, d.MaxBodyBytes)
	rl := middleware.NewRateLimiter(d.RateLimitPerMinute)
	r.With(rl.Middleware).HandleFunc("/hook/{binID}", ingestH.Ingest)

	binH := NewBinHandler(d.Bins)
	reqH := NewRequestHandler(d.Requests)
	replayH := NewReplayHandler(d.Replay)

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.Auth(d.AdminAPIKey))

		r.Post("/bins", binH.Create)
		r.Get("/bins", binH.List)
		r.Get("/bins/{id}", binH.Get)
		r.Delete("/bins/{id}", binH.Delete)

		r.Get("/bins/{id}/requests", reqH.ListByBin)
		r.Get("/requests/{id}", reqH.Get)
		r.Delete("/requests/{id}", reqH.Delete)
		r.Post("/requests/{id}/replay", replayH.Replay)
	})

	return r
}
