package web

import (
	"net/http"
)

/*
An HTTP multiplexer, much like net/http's ServeMux.

Routes may be added using any of the various HTTP-method-specific functions.
When processing a request, when iterating in insertion order the first route
that matches both the request's path and method is used.

There are two other differences worth mentioning between web.Mux and
http.ServeMux. First, string patterns (i.e., Sinatra-like patterns) must match
exactly: the "rooted subtree" behavior of ServeMux is not implemented. Secondly,
unlike ServeMux, Mux does not support Host-specific patterns.

If you require any of these features, remember that you are free to mix and
match muxes at any part of the stack.

In order to provide a sane API, many functions on Mux take interface{}'s. This
is obviously not a very satisfying solution, but it's probably the best we can
do for now. Instead of duplicating documentation on each method, the types
accepted by those functions are documented here.

A middleware (the untyped parameter in Use() and Insert()) must be one of the
following types:
	- func(http.Handler) http.Handler
	- func(c *web.C, http.Handler) http.Handler
All of the route-adding functions on Mux take two untyped parameters: pattern
and handler. Pattern must be one of the following types:
	- string (interpreted as a Sinatra pattern)
	- regexp.Regexp
	- web.Pattern
Handler must be one of the following types:
	- http.Handler
	- web.Handler
	- func(w http.ResponseWriter, r *http.Request)
	- func(c web.C, w http.ResponseWriter, r *http.Request)
*/
type Mux struct {
	mStack
	router
}

// Create a new Mux without any routes or middleware.
func New() *Mux {
	mux := Mux{
		mStack: mStack{
			stack: make([]mLayer, 0),
			pool:  make(chan *cStack, mPoolSize),
		},
		router: router{
			routes:   make([]route, 0),
			notFound: parseHandler(http.NotFound),
		},
	}
	mux.mStack.router = HandlerFunc(mux.router.route)
	return &mux
}

// Serve a request with the given Mux. Satisfies the http.Handler interface.
func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	stack := m.mStack.alloc()
	defer m.mStack.release(stack)

	stack.ServeHTTP(w, r)
}

// Serve a context dependent request with the given Mux. Satisfies the
// web.Handler interface.
func (m *Mux) ServeHTTPC(c C, w http.ResponseWriter, r *http.Request) {
	stack := m.mStack.alloc()
	defer m.mStack.release(stack)

	stack.ServeHTTPC(c, w, r)
}