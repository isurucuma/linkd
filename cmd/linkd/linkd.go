package main

import (
	"errors"
	"linkd/httpio"
	"linkd/link"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	const timeout = 10 * time.Second
	const addr = "localhost:8080"

	// log := slog.With("app", "linkd")
	log := slog.New(&httpio.LogHandler{
		Handler: slog.NewTextHandler(os.Stderr, nil),
	})

	log = log.With("app", "linkd")

	slog.SetDefault(log)

	log.Info("starting", "addr", addr)

	links := link.NewServer(link.NewStore())
	handler := http.TimeoutHandler(links, time.Second, "timeout")
	handler = httpio.LoggingMiddleware(handler)
	handler = httpio.TraceMiddleware(handler)

	srv := &http.Server{
		Addr:        addr,
		ReadTimeout: timeout * 2,
		IdleTimeout: timeout * 4,
		Handler:     handler,
	}

	err := srv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) { // ErrServerClosed is a expedted if not that is an unexpected error
		log.Error("server closed unexpectedly", "message", err)
	}
}
