package chi

import (
	"net/http"

	"final/internal/ports"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// ChiRouter реализует ports.Router
type ChiRouter struct {
	mux *chi.Mux
}

func NewChiRouter() *ChiRouter {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.RealIP)
	return &ChiRouter{mux: mux}
}

// Get регистрирует GET маршрут
func (r *ChiRouter) Get(pattern string, handler http.HandlerFunc) {
	r.mux.Get(pattern, handler)
}

// Post регистрирует POST маршрут
func (r *ChiRouter) Post(pattern string, handler http.HandlerFunc) {
	r.mux.Post(pattern, handler)
}

// Put регистрирует PUT маршрут
func (r *ChiRouter) Put(pattern string, handler http.HandlerFunc) {
	r.mux.Put(pattern, handler)
}

// Delete регистрирует DELETE маршрут
func (r *ChiRouter) Delete(pattern string, handler http.HandlerFunc) {
	r.mux.Delete(pattern, handler)
}

// Patch регистрирует PATCH маршрут
func (r *ChiRouter) Patch(pattern string, handler http.HandlerFunc) {
	r.mux.Patch(pattern, handler)
}

// Use добавляет middleware к текущему роутеру
func (r *ChiRouter) Use(middlewares ...func(http.Handler) http.Handler) {
	r.mux.Use(middlewares...)
}

// With создаёт новый роутер с дополнительными middleware
func (r *ChiRouter) With(middlewares ...func(http.Handler) http.Handler) ports.Router {
	return &ChiRouter{
		mux: r.mux.With(middlewares...).(*chi.Mux),
	}
}

// Group создаёт вложенную группу маршрутов
func (r *ChiRouter) Group(fn func(r ports.Router)) ports.Router {
	r.mux.Group(func(g chi.Router) {
		fn(&ChiRouter{mux: g.(*chi.Mux)})
	})
	return r
}

// ServeHTTP реализует http.Handler
func (r *ChiRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
