package httpadapter

import "net/http"

// Router — контракт HTTP-адаптера; не является портом прикладного слоя.
type Router interface {
	http.Handler
	Get(pattern string, handler http.HandlerFunc)
}
