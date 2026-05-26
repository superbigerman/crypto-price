package ports

import "net/http"

type Router interface {
	http.Handler
	Get(pattern string, handler http.HandlerFunc)
	Post(pattern string, handler http.HandlerFunc)
	Put(pattern string, handler http.HandlerFunc)
	Delete(pattern string, handler http.HandlerFunc)
	Patch(pattern string, handler http.HandlerFunc)
	Use(middlewares ...func(http.Handler) http.Handler)

	With(middlewares ...func(http.Handler) http.Handler) Router
	Group(fn func(r Router)) Router
}
