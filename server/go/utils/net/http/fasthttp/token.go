package fasthttpi

import (
	"net/url"

	"github.com/liov/hoper/v2/utils/net/http"
	stringsi "github.com/liov/hoper/v2/utils/strings"
	"github.com/valyala/fasthttp"
)

func GetToken(req *fasthttp.Request) string {
	if token := stringsi.ToString(req.Header.Peek(httpi.HeaderAuthorization)); token != "" {
		return token
	}
	if cookie := stringsi.ToString(req.Header.Cookie("token")); len(cookie) > 0 {
		token, _ := url.QueryUnescape(cookie)
		return token
	}
	return ""
}
