package fasthttpi

import (
	"net/url"

	stringsi "github.com/liov/hoper/go/v2/utils/strings"
	"github.com/valyala/fasthttp"
)

func GetToken(req *fasthttp.Request) string {
	var token string
	cookie := stringsi.ToString(req.Header.Cookie("token"))
	if len(cookie) > 0 {
		token, _ = url.QueryUnescape(cookie)
		return token
	}
	return stringsi.ToString(req.Header.Peek("Authorization"))
}
