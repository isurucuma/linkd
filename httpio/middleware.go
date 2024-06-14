package httpio

import (
	"log/slog"
	"net/http"
	"time"
)

type ResponseRecorder struct {
	http.ResponseWriter
	status int
}

func (r *ResponseRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rr := &ResponseRecorder{
			ResponseWriter: w,
		}
		next.ServeHTTP(rr, r)
		slog.Log(r.Context(), slog.LevelInfo, "request",
			"url", r.URL,
			"method", r.Method,
			"status", rr.status,
			"took", time.Since(start))
	}
}
