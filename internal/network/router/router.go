package router

import (
	"io"
	"net/http"
	"net/http/httptest"
)


type Router struct {
	mux *http.ServeMux
}

func NewRouter() *Router {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
		w.WriteHeader(http.StatusOK)

		w.Header().Set("AtEnd1", "value 1")
		io.WriteString(w, "This HTTP response has both headers before this text and trailers at the end.\n")
	})



	return &Router{
		mux: mux,
	}
}

func (r *Router) GetResponse(req *http.Request) *http.Response {
	h, _ := r.mux.Handler(req)
	rw := httptest.NewRecorder()

	h.ServeHTTP(rw, req)
	return rw.Result()
}
