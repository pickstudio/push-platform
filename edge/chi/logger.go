package chi

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

func AccessLoggerHandler(logger zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			before := time.Now()
			next.ServeHTTP(ww, r)

			scheme := "http"
			if r.TLS != nil {
				scheme = "https"
			}
			logger.Info().
				Str("request_id", middleware.GetReqID(r.Context())).
				Str("url", fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)).
				Str("uri", r.RequestURI).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("query", r.URL.RawQuery).
				Str("remote_ip", r.RemoteAddr).
				Str("proto", r.Proto).
				Str("schema", scheme).
				Int("status", ww.Status()).
				Float64("duration", float64(time.Since(before))/1000000.0).
				Send()
		}
		return http.HandlerFunc(fn)
	}
}
