package chi

import (
	"net/http"

	chiprometheus "github.com/daangn/go-chi-prometheus"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	corsmiddleware "github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"

	_const "github.com/pickstudio/push-platform/const"
)

func New() chi.Router {
	r := chi.NewRouter()
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.Recoverer)
	r.Use(AccessLoggerHandler(
		log.With().
			Str("service", _const.Project).
			Str(_const.KeyLogType, _const.ValueAccessLog).
			Logger(),
	))
	r.Use(corsmiddleware.New(corsmiddleware.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           30,
	}).Handler)
	r.Use(chiprometheus.NewMiddleware(chiprometheus.WithHistogram(true)))
	r.Use(chimiddleware.RequestID)
	r.Handle("/metrics", promhttp.Handler())
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	return r
}
