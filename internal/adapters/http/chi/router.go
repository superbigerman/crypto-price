package chiadapter

import (
	"net/http"

	httpadapter "final/internal/adapters/http"

	"github.com/go-chi/chi/v5"
)

var _ httpadapter.Router = (*Router)(nil)

type Router struct {
	mux *chi.Mux
}

func NewRouter() *Router {
	return &Router{mux: chi.NewRouter()}
}

func (r *Router) Get(pattern string, handler http.HandlerFunc) {
	r.mux.Get(pattern, handler)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
