package httpi

import (
	"net/http"
	"net/url"
)

func GetToken(r *http.Request) string {
	cookie, _ := r.Cookie("token")
	if cookie != nil {
		value, _ := url.QueryUnescape(cookie.Value)
		return value
	}
	return r.Header.Get("authorization")
}
