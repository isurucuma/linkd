package link

import (
	"context"
	"linkd/sqlx/sqlxtest"
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

func TestShorten(t *testing.T) {
	t.Parallel()
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(
		`{"key": "go", "url": "https://go.dev"}`,
	))
	store := NewStore(
		sqlxtest.Dial(context.Background(), t),
	)
	Shorten(store).ServeHTTP(w, r)
	if w.Code != http.StatusCreated {
		t.Errorf("got status code = %d, want %d", w.Code, http.StatusCreated)
	}
	if got := w.Body.String(); !strings.Contains(got, `"go"`) {
		t.Errorf("got body = %s\twant contains %s", got, `"go"`)
	}
}

type fakeLinkCreator func(context.Context, Link) error

func (f fakeLinkCreator) Create(ctx context.Context, link Link) error {
	return f(ctx, link)
}

func TestShortenWithTestDouble(t *testing.T) {
	creator := fakeLinkCreator(func(ctx context.Context, l Link) error {
		return nil
	})
	handler := shorten(creator)

	body := `{"key": "go", "url": "https://go.dev"}`
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(body))
	handler.ServeHTTP(w, r)

	if want := http.StatusCreated; w.Code != want {
		t.Errorf("got status code = %d, want %d", w.Code, want)
	}
}

func TestShortenInternalError(t *testing.T) {
	t.Parallel()
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(
		`{"key": "go", "url": "https://go.dev"}`,
	))
	store := NewStore(
		sqlxtest.Dial(context.Background(), t))
	_ = store.db.Close()
	Shorten(store).ServeHTTP(w, r)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("got status code = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}
