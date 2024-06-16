package link

import (
	"fmt"
	"linkd/bite"
	"linkd/httpio"
	"net/http"
)

type Server struct {
	http.Handler
}

func NewServer(links *Store) *Server {
	mux := http.NewServeMux()
	mux.Handle("POST /shorten", Shorten(links))
	mux.Handle("GET /r/{key}", Resolve(links))
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

// Shorten handles the URL shortening requests.
//
//	Status Code       Condition
//	201               The link is successfully shortened.
//	400               The request is invalid.
//	409               The link already exists.
//	405               The request method is not POST.
//	413               The request body is too large.
//	500               There is an internal error.
func Shorten(links *Store) httpio.Handler {
	return func(w http.ResponseWriter, r *http.Request) httpio.Handler {
		// link := Link{
		// 	Key: r.FormValue("key"),
		// 	URL: r.FormValue("url"),
		// }
		var link Link
		max := http.MaxBytesReader(w, r.Body, 4_096) // provides extra layer of protection, if the body is larger than the max bytes then this will close the reader
		if err := httpio.DecodeJSON(max, &link); err != nil {
			return httpio.Error(invalidRequest(err))
		}
		if err := links.Create(r.Context(), link); err != nil {
			// httpio.Error(err)
			// return
			return httpio.Error(err)
		}
		// w.WriteHeader(http.StatusCreated)
		// w.Write([]byte(link.Key))
		return httpio.Code(http.StatusCreated, httpio.JSON(
			map[string]string{
				"key": link.Key,
			},
		))
	}
}

// Resolve handles the URL resolving requests for the short links.
//
//	Status Code       Condition
//	302               The link is successfully resolved.
//	400               The request is invalid.
//	404               The link does not exist.
//	500               There is an internal error.
func Resolve(links *Store) httpio.Handler {
	return func(w http.ResponseWriter, r *http.Request) httpio.Handler {
		link, err := links.Retrieve(r.Context(), r.PathValue("key"))
		if err != nil {
			return httpio.Error(err)
		}
		http.Redirect(w, r, link.URL, http.StatusFound)
		return httpio.Ok
	}
}

func Health(w http.ResponseWriter, r *http.Request) { // by default handlers write a status code of OK
	fmt.Fprintln(w, "ok")
}

func invalidRequest(err error) error {
	return fmt.Errorf("%w: %v", bite.ErrInvalidRequest, err)
}
