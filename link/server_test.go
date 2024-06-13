package link

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealth(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	Health(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("got status code %d, want %d", w.Code, http.StatusOK)
	}
	if got := w.Body.String(); !strings.Contains(got, "ok") {
		t.Errorf("got body = %s\twant contains %s", got, "OK")
	}
}

func TestServer(t *testing.T) {
	t.Parallel()

	testData := []struct {
		path               string
		method             string
		expectedStatusCode int
	}{
		{path: "/health", method: http.MethodGet, expectedStatusCode: http.StatusOK},
		{path: "/notFound", method: http.MethodGet, expectedStatusCode: http.StatusNotFound},
	}

	for _, tt := range testData {
		t.Run(tt.path, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(tt.method, tt.path, nil)

			srv := NewServer(nil)
			srv.ServeHTTP(w, r)

			if w.Code != tt.expectedStatusCode {
				t.Errorf("got status code = %d, want %d", w.Code, tt.expectedStatusCode)
			}
		})
	}
}
