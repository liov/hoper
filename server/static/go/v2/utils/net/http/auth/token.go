package auth

import (
	"net/http"
	"net/url"
)

func GetToken(r *http.Request) string {
	cookie, _ := r.Cookie("token")
	value, _ := url.QueryUnescape(cookie.Value)
	if value != "" {
		return value
	}
	return r.Header.Get("authorization")
}

