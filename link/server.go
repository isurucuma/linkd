package link

import (
	"errors"
	"fmt"
	"linkd/bite"
	"net/http"
)

type Server struct {
	http.Handler
}

func NewServer(links *Store) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /shorten", Shorten(links))
	mux.HandleFunc("GET /r/{key}", Resolve(links))
	mux.HandleFunc("GET /health", Health)
	return &Server{
		Handler: mux,
	}
}

// func Shorten(w http.ResponseWriter, r *http.Request) {
// 	w.WriteHeader(http.StatusCreated)
// 	fmt.Fprintln(w, "go")
// }

// func Resolve(w http.ResponseWriter, r *http.Request) {
// 	key := r.PathValue("key")
// 	fmt.Println("Key is: ", key)
// 	const uri = "https://go.dev"
// 	http.Redirect(w, r, uri, http.StatusFound)
// }

func Shorten(links *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		link := Link{
			Key: r.FormValue("key"),
			URL: r.FormValue("url"),
		}
		if err := links.Create(r.Context(), link); err != nil {
			httpError(w, err)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(link.Key))
	}
}

func Resolve(links *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		link, err := links.Retrieve(r.Context(), r.PathValue("key"))
		if err != nil {
			httpError(w, err)
			return
		}
		http.Redirect(w, r, link.URL, http.StatusFound)
	}
}

func Health(w http.ResponseWriter, r *http.Request) { // by default handlers write a status code of OK
	fmt.Fprintln(w, "ok")
}

func httpError(w http.ResponseWriter, err error) {
	if err == nil { // no error #A
		return
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
	http.Error(w, err.Error(), code)
}
