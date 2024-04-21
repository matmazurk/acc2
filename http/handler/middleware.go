package handler

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

func logh(next http.Handler, logger zerolog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		cw := &statusCaptureResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		next.ServeHTTP(cw, r)

		logger.Info().Str("path", r.URL.Path).Str("duration", time.Since(start).String()).Int("response_code", cw.statusCode).Msg("request handled")
	})
}

type statusCaptureResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusCaptureResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
