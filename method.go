package muxify

import "net/http"

func (m *Mux) Use(middleware func(http.Handler) http.Handler) {
	m.middlewares = append(m.middlewares, middleware)

}

func (m *Mux) Get(path string, h http.HandlerFunc) {
	m.Handle(http.MethodGet, path, h)
}

func (m *Mux) Post(path string, h http.HandlerFunc) {
	m.Handle(http.MethodPost, path, h)
}

func (m *Mux) Put(path string, h http.HandlerFunc) {
	m.Handle(http.MethodPut, path, h)
}

func (m *Mux) Delete(path string, h http.HandlerFunc) {
	m.Handle(http.MethodDelete, path, h)
}
