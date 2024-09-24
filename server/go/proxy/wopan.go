package main

import (
	"github.com/hopeio/utils/log"
	"github.com/rs/cors"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	proxy := &httputil.ReverseProxy{Rewrite: func(r *httputil.ProxyRequest) {
		targets := r.In.Header["Target-Url"]
		if len(targets) == 0 {
			log.Error("no target specified")
			return
		}
		target := targets[0]
		targetUrl, _ := url.Parse(target)
		r.Out.URL = r.In.URL
		r.Out.Host = targetUrl.Host
		r.Out.URL.Host = targetUrl.Host
		r.Out.URL.Scheme = "https"

		/*		r.Out.Header["Refer"] = r.In.Header["Target-Refer"]
				r.Out.Header["Origin"] = r.In.Header["Target-Origin"]*/
	},
		Transport: &http.Transport{
			Proxy:             http.ProxyFromEnvironment, // 代理使用
			ForceAttemptHTTP2: true,
		}}
	server := cors.AllowAll()

	log.Fatal(http.ListenAndServe(":8080", server.Handler(proxy)))
}
