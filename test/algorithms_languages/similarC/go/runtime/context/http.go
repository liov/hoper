package main

import (
	"context"
	"net/http"
	"time"
)

// ContextMiddle是http服务中间件，统一读取通行cookie并使用ctx传递
func ContextMiddle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, _ := r.Cookie("Check")
		if cookie != nil {
			ctx := context.WithValue(r.Context(), "Check", cookie.Value)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

// 强制设置通行cookie
func CheckHandler(w http.ResponseWriter, r *http.Request) {
	expitation := time.Now().Add(24 * time.Hour)
	cookie := http.Cookie{Name: "Check", Value: "42", Expires: expitation}
	http.SetCookie(w, &cookie)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// 通过取中间件传过来的context值来判断是否放行通过
	if chk := r.Context().Value("Check"); chk == "42" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Let's go! \n"))
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No Pass!"))
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler)

	// 人为设置通行cookie
	mux.HandleFunc("/chk", CheckHandler)

	ctxMux := ContextMiddle(mux)
	http.ListenAndServe(":8080", ctxMux)
}
