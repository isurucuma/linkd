package httpio

import (
	"errors"
	"fmt"
	"linkd/bite"
	"log/slog"
	"net/http"
)

type Handler func(w http.ResponseWriter, r *http.Request) Handler

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if next := h(w, r); next != nil {
		next.ServeHTTP(w, r)
	}
}

func Ok(w http.ResponseWriter, r *http.Request) Handler {
	return nil
}

func Code(statusCode int, next Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) Handler {
		w.WriteHeader(statusCode)
		return next
	}
}

func Text(s string) Handler {
	return func(w http.ResponseWriter, r *http.Request) Handler {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintln(w, s)
		return Ok
	}
}

func CodeText(statusCode int, msg string) Handler {
	return func(w http.ResponseWriter, r *http.Request) Handler {
		return Code(statusCode, Text(msg))
	}
}

func Error(err error) Handler {
	if err == nil { // no error
		return Ok
	}
	var code int
	switch {
	case errors.Is(err, bite.ErrInvalidRequest):
		code = http.StatusBadRequest
	case errors.Is(err, bite.ErrExists):
		code = http.StatusConflict
	case errors.Is(err, bite.ErrNotExists):
		code = http.StatusNotFound
	default:
		code = http.StatusInternalServerError
	}

	return func(w http.ResponseWriter, r *http.Request) Handler {
		if code == http.StatusInternalServerError {
			slog.Log(r.Context(), slog.LevelError, "internal", "url", r.URL, "message", err)
			err = bite.ErrInternal
		}
		return Code(code, Text(err.Error()))
	}

}
