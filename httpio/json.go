package httpio

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

func JSON(v any) Handler {
	return func(w http.ResponseWriter, r *http.Request) Handler {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err := json.NewEncoder(w).Encode(v); err != nil {
			slog.Log(r.Context(), slog.LevelError, "internal", "url", r.URL, "message", err)
		}
		return Ok
	}
}

func DecodeJSON(r io.Reader, v any) error {
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()
	return decoder.Decode(v)
}
