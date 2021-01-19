package fasthttpi

import (
	"net/url"

	"github.com/liov/hoper/go/v2/utils/net/http/request"
	stringsi "github.com/liov/hoper/go/v2/utils/strings"
	"github.com/valyala/fasthttp"
)

func GetToken(req *fasthttp.Request) string {
	if token:=stringsi.ToString(req.Header.Peek(request.Authorization));token!=""{
		return token
	}
	if cookie := stringsi.ToString(req.Header.Cookie("token"));len(cookie) > 0 {
		token, _ := url.QueryUnescape(cookie)
		return token
	}
	return ""
}
