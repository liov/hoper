package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// FromStd converts native http.Handler & http.HandlerFunc to gin.HandlerFunc.
//
// Supported form types:
// 		 .FromStd(h http.Handler)
// 		 .FromStd(func(w http.ResponseWriter, r *http.Request))
// 		 .FromStd(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc))
func FromStd(handler interface{}) gin.HandlerFunc {
	switch h := handler.(type) {
	case gin.HandlerFunc:
		return h
	case http.Handler:
		return func(ctx *gin.Context) {
			h.ServeHTTP(ctx.Writer, ctx.Request)
		}
	case func(http.ResponseWriter, *http.Request):
		return FromStd(http.HandlerFunc(h))
	case func(http.ResponseWriter, *http.Request, http.HandlerFunc):
		return FromStdWithNext(h)

	default:
		{
			panic(fmt.Errorf(`
			Passed argument is not a func(*gin.Context) neither one of these types:
			- http.Handler
			- func(w http.ResponseWriter, r *http.Request)
			- func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
			---------------------------------------------------------------------
			It seems to be a %T points to: %v`, handler, handler))
		}
	}
}

func FromStdWithNext(h func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		next := func(w http.ResponseWriter, r *http.Request) {
			ctx.Request = r
			ctx.Next()
		}

		h(ctx.Writer, ctx.Request, next)
	}
}

func Converts(handlers []http.HandlerFunc) []gin.HandlerFunc {
	var rets []gin.HandlerFunc
	for _, handler := range handlers {
		rets = append(rets, FromStd(handler))
	}
	return rets
}

func Convert(handler http.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		handler(ctx.Writer, ctx.Request)
	}
}
