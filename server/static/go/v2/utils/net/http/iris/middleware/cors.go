package middleware

import (
	"github.com/kataras/iris/v12"
)

func Cors(ctx iris.Context) {
	/*		//method := c.Request.Method

			origin := c.Request.Header.Get("Origin")
			var headerKeys []string
			for k := range c.Request.Header {
				headerKeys = append(headerKeys, k)
			}
			headerStr := strings.Join(headerKeys, ", ")
			if headerStr != "" {
				headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
			} else {
				headerStr = "access-control-allow-origin, access-control-allow-headers"
			}
			if origin != "" {
				//下面的都是乱添加的-_-~

				c.Header("Access-Control-Allow-Origin", "*")
				c.Header("Access-Control-Allow-Headers", headerStr)
				c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
				c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
				c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
				// c.Header("Access-Control-Max-Age", "172800")
				c.Header("Access-Control-Allow-Credentials", "true")
				c.Set("content-type", "application/json")
			}

			//放行所有OPTIONS方法
			if c.Request.Method == "OPTIONS" {
				c.JSON(http.StatusOK, "Options Request!")
			}*/

	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Headers", "Accept、Accept-Language,Content-Language,Content-Type")
	ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
	// c.Header("Access-Control-Max-Age", "172800")
	ctx.Header("Access-Control-Allow-Credentials", "true")
	ctx.ContentType("application/json")
	if ctx.Method() == "OPTIONS" {
		ctx.JSON("Options Request!")
	}
}
