package httpio

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoggingMiddleware(t *testing.T) {
	var buf bytes.Buffer
	log := slog.New(slog.NewTextHandler(&buf, nil))
	slog.SetDefault(log)

	var handler http.Handler
	handler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(418)
	})

	handler = LoggingMiddleware(handler)
	handler.ServeHTTP(
		httptest.NewRecorder(),
		httptest.NewRequest("GET", "/test", nil),
	)

	got := buf.String()
	fmt.Println(got)
	if !strings.Contains(got, "GET") {
		t.Error("want GET in the log")
	}
	if !strings.Contains(got, "418") {
		t.Error("want 418 in the log")
	}
	if !strings.Contains(got, "/test") {
		t.Errorf("want /test in the log")
	}
	if t.Failed() {
		t.Log("got:", got)
	}
}

func TestTraceMiddleware(t *testing.T) {
	var (
		traceId         int64
		traceValPresent bool
	)

	var handler http.Handler
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceId, traceValPresent = TraceId(r.Context())
	})

	handler = TraceMiddleware(handler)

	handler.ServeHTTP(
		httptest.NewRecorder(),
		httptest.NewRequest("GET", "/test", nil),
	)

	if !traceValPresent {
		t.Fatalf("got context without trace id")
	}

	if traceId <= 0 {
		t.Fatalf("got %d, want a positive trave id", traceId)
	}

	prevTraceId := traceId
	handler.ServeHTTP(
		httptest.NewRecorder(),
		httptest.NewRequest("GET", "/test", nil),
	)

	if prevTraceId == traceId {
		t.Fatalf("got duplicate trace id: %d", traceId)
	}
}

func TestLogHandler(t *testing.T) {
	var buf bytes.Buffer

	ctx := SetTraceId(context.Background(), 22)
	testLogger := slog.New(&LogHandler{
		Handler: slog.NewTextHandler(&buf, nil),
	})

	testLogger.Log(ctx, slog.LevelInfo, "test")

	if got := buf.String(); !strings.Contains(got, "22") {
		t.Errorf("want trace id %d in the log, got %s", 22, got)
	}
}
