package handler

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

func logh(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		cw := &statusCaptureResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		next.ServeHTTP(cw, r)

		log.Info().Str("path", r.URL.Path).Str("duration", time.Since(start).String()).Int("response_code", cw.statusCode).Msg("request handled")
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
