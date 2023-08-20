// go_mux.go

package go_mux

import "net/http"

type Mux struct {
	mux *http.ServeMux
}

func NewMux() *Mux {
	return &Mux{
		mux: http.NewServeMux(),
	}
}

func (cm *Mux) Handle(pattern string, handler http.Handler) {
	cm.mux.Handle(pattern, handler)
}

func (cm *Mux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	cm.mux.HandleFunc(pattern, handler)
}

func (cm *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cm.mux.ServeHTTP(w, r)
}
