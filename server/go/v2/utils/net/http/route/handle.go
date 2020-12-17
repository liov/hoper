package route

import "net/http"

type Handlers []http.Handler

func (hs Handlers) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, handler := range hs {
		handler.ServeHTTP(w, req)
	}
}

type HandlerFuncs []http.HandlerFunc

func (hs HandlerFuncs) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, handler := range hs {
		handler(w, req)
	}
}

func (hs HandlerFuncs) HandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		for _, handler := range hs {
			handler(w, req)
		}
	}
}

func (hs *HandlerFuncs) Add(handler http.HandlerFunc) {
	*hs = append(*hs, handler)
}
