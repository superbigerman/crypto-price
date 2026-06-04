package chi

import (
	"net/http"

	"final/internal/ports"

	"github.com/go-chi/chi/v5"
)

var _ ports.Router = (*ChiRouter)(nil)

type ChiRouter struct {
	mux *chi.Mux //
}

func NewChiRouter() *ChiRouter {
	return &ChiRouter{
		mux: chi.NewRouter(),
	}
}

func (r *ChiRouter) Get(pattern string, handler http.HandlerFunc) {
	r.mux.Get(pattern, handler)
}

func (r *ChiRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
