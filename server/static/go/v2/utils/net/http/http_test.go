package http

import (
	"net/http"
	"testing"

	"github.com/kataras/iris/v12"
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
	//g := gin.New()
	//g.GET("/:id", func(context *gin.Context) {})
	//g.GET("/:name", func(context *gin.Context) {})
	//g.GET("/*file", func(context *gin.Context) {})

	i := iris.New()
	i.Get("/:id/:name/:path", func(context iris.Context) { context.WriteString("/:id/:name/:path") })
	i.Get("/id/name/path", func(context iris.Context) { context.WriteString("/id/name/path") })
	i.Run(iris.Addr(":8080"))
}
