package httpi

import (
	"net/http"
	"net/url"

	"github.com/liov/hoper/go/v2/utils/net/http/request"
)

func GetToken(r *http.Request) string {
	if token := r.Header.Get(request.Authorization); token != "" {
		return token
	}
	if cookie, _ := r.Cookie(request.Token); cookie != nil {
		value, _ := url.QueryUnescape(cookie.Value)
		return value
	}
	return ""
}
