package httpio

import (
	"bytes"
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
