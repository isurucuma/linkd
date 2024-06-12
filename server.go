package link

import (
	"fmt"
	"net/http"
)

type Server struct {
	http.Handler
}

func NewServer() *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /shorten", Shorten)
	mux.HandleFunc("GET /r/{key}", Resolve)
	mux.HandleFunc("GET /health", Health)
	return &Server{
		Handler: mux,
	}
}

func Shorten(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "go")
}

func Resolve(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	fmt.Println("Key is: ", key)
	const uri = "https://go.dev"
	http.Redirect(w, r, uri, http.StatusFound)
}

func Health(w http.ResponseWriter, r *http.Request) { // by default handlers write a status code of OK
	fmt.Fprintln(w, "ok")
}
