package main

import "net/http"

func main() {

}

//无语，返回值其实是func(http.Handler) http.Handler
func CORS() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return h
	}
}

type Debug interface {
	Debug() string
}

type DbgFunc func(debug Debug) Debug

func (d DbgFunc) Debug() string {
	return "d()"
}

/*func ImpDbg() Debug {
	return func(d Debug) Debug {
		return d
	}
}*/
