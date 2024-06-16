package httpio

import (
	"context"
	"log/slog"
	"math/rand"
	"net/http"
	"time"
)

type ResponseRecorder struct {
	http.ResponseWriter
	status int
}

type traceIdKey struct{}

type LogHandler struct{ slog.Handler }

func (r *ResponseRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *ResponseRecorder) Unwrap() http.ResponseWriter {
	return r.ResponseWriter
}

func LoggingMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rr := &ResponseRecorder{
			ResponseWriter: w,
		}
		// rc := http.NewResponseController(w)
		// the type which is passed as the ResponseWritter is not only a response writter, it will have the implementations of some other methods
		// of some interfaces like Flusher. Therefore if we need to call those methods then we must be able to
		// making a ResponseController from that ResponseWritter is the way to access those methods
		// ResponseController is a struct which is having all the underline methods that the ResponseWritter is not defined
		// then we can call those methods using the ResponseController. For that to happen we need to have a Unwrap method in the ResponseRecoder(what ever the type that you created)
		// that Unwrap method should return the underline ResponseWrtiter
		// if err := rc.Flush(); errors.Is(err, http.ErrNotSupported) {
		//   /* sorry, ResponseWriter does not support Flusher */
		// }

		next.ServeHTTP(rr, r)
		slog.Log(r.Context(), slog.LevelInfo, "request",
			"url", r.URL,
			"method", r.Method,
			"status", rr.status,
			"took", time.Since(start))
	}
}

func TraceMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := SetTraceId(r.Context(), rand.Int63())
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func SetTraceId(ctx context.Context, id int64) context.Context {
	return context.WithValue(ctx, traceIdKey{}, id)
}

func TraceId(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(traceIdKey{}).(int64)
	return id, ok
}

func (h *LogHandler) Handle(ctx context.Context, r slog.Record) error {
	if id, ok := TraceId(ctx); ok {
		r.Add("trace_id", id)
	}
	return h.Handler.Handle(ctx, r)
}

func (h *LogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &LogHandler{Handler: h.Handler.WithAttrs(attrs)}
}

func (h *LogHandler) WithGroup(name string) slog.Handler {
	return &LogHandler{Handler: h.Handler.WithGroup(name)}
}
