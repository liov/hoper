package http

import (
	"net/http"
	"testing"
)

func TestRoute(t *testing.T) {
	sv := http.DefaultServeMux
	var f = func(writer http.ResponseWriter, request *http.Request) {}
	sv.HandleFunc("/", f)
	sv.HandleFunc("/a", f)
	sv.HandleFunc("/a/", f)
	sv.HandleFunc("/a/b", f)
	sv.HandleFunc("/b", f)
	sv.HandleFunc("/b/a/c/d", f)
	sv.HandleFunc("/c/", f)
}
