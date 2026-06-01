package ports

import "net/http"

type Router interface {
	http.Handler
	Get(pattern string, handler http.HandlerFunc)
}
