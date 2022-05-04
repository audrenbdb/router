package router

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type router struct {
	muxMutex  sync.Mutex
	mux       *http.ServeMux
	endpoints []endpoint
}

// New creates a router based on http.ServeMux.
//
// Router has the ability to register an endpoint handler
// from one method, such as:
//
// r.GET("/foo", handleBar)
func New() *router {
	return &router{
		endpoints: make([]endpoint, 0),
	}
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if r.mux == nil {
		r.initializeMux()
	}
	r.mux.ServeHTTP(w, req)
}

func (r *router) initializeMux() {
	r.muxMutex.Lock()
	defer r.muxMutex.Unlock()
	if r.mux != nil {
		return
	}
	r.mux = http.NewServeMux()
	for _, ep := range r.endpoints {
		r.mux.Handle(ep.pattern, ep.handler)
	}
}

func (r *router) GET(pattern string, handlerFunc http.HandlerFunc) {
	r.registerEndpoint(http.MethodGet, pattern, handlerFunc)
}

func (r *router) POST(pattern string, handlerFunc http.HandlerFunc) {
	r.registerEndpoint(http.MethodPost, pattern, handlerFunc)
}

func (r *router) PUT(pattern string, handlerFunc http.HandlerFunc) {
	r.registerEndpoint(http.MethodPut, pattern, handlerFunc)
}

func (r *router) PATCH(pattern string, handlerFunc http.HandlerFunc) {
	r.registerEndpoint(http.MethodPatch, pattern, handlerFunc)
}

func (r *router) DELETE(pattern string, handlerFunc http.HandlerFunc) {
	r.registerEndpoint(http.MethodDelete, pattern, handlerFunc)
}

type endpoint struct {
	pattern string
	handler methodHandler
}

// methodHandler is a map of type map[request method]http.Handler
//
// it implements http.Handler
type methodHandler map[string]http.Handler

func (m methodHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h, ok := m[r.Method]
	if !ok {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	h.ServeHTTP(w, r)
}

// registerEndpoint adds an endpoint handler
// from given method to router.
//
// an endpoint with same method cannot be registered twice
func (r *router) registerEndpoint(method, pattern string, handler http.Handler) {
	ep := r.findEndpoint(pattern)
	if ep == nil {
		newEndpoint := endpoint{
			pattern: pattern,
			handler: map[string]http.Handler{},
		}
		r.endpoints = append(r.endpoints, newEndpoint)
		ep = &newEndpoint
	}
	_, methodFound := ep.handler[method]
	if methodFound {
		msg := fmt.Sprintf("method %s already registered for pattern %s", method, pattern)
		log.Fatal(msg)
	}
	ep.handler[method] = handler
}

func (r *router) findEndpoint(pattern string) *endpoint {
	for _, ep := range r.endpoints {
		if ep.pattern == pattern {
			return &ep
		}
	}
	return nil
}
