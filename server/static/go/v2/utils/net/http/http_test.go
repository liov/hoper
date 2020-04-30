package http

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
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
	g := gin.New()
	g.GET("/:id", func(context *gin.Context) {})
	//g.GET("/:name", func(context *gin.Context) {})
	g.GET("/*file", func(context *gin.Context) {})
}
