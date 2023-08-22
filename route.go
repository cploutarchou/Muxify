package muxify

import "net/http"

func NewRouter() *Mux {
	return NewMux()
}

type Router interface {
	Handle(method, path string, handler http.Handler)
	http.Handler
	Routes
	Use(middlewares ...func(http.Handler) http.Handler)
	Route(pattern string, fn func(r Router)) Router
}

type Context struct {
	Response *http.ResponseWriter
	Request  *http.Request
	Params   map[string]string
}

type Routes interface {
	Routes() []Route
	Middlewares() Middlewares
	Match(ctx *Context, method, path string) bool
}

type Middlewares []func(http.Handler) http.Handler

type Route struct {
	Method          string
	Path            string
	Handler         http.Handler
	Middlewares     Middlewares
	Name            string
	SubRoutes       []Route
	SubRouters      []Router
	SubMuxes        []*Mux
	SubHandlers     []http.Handler
	HandlerFuncs    []func(http.ResponseWriter, *http.Request)
	HandlerFuncsCtx []func(*Context)
}

func (r *Route) Routes() []Route {
	return r.SubRoutes
}
