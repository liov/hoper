package client

import (
	"context"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"io"
	"net/http"
	urlpkg "net/url"
	"strings"
)

type replaceHttpRequest http.Request

func NewReplaceHttpRequest(r *http.Request) *replaceHttpRequest {
	return (*replaceHttpRequest)(r)
}

func ReplaceHttpRequest(r *http.Request, url, method string, body io.ReadCloser) *http.Request {
	return NewReplaceHttpRequest(r).Replace(url, method, body).StdHttpRequest()
}

func (r *replaceHttpRequest) Replace(url, method string, body io.ReadCloser) *replaceHttpRequest {
	r.SetURL(url).SetMethod(method).SetBody(body)
	return r
}

func (r *replaceHttpRequest) SetURL(url string) *replaceHttpRequest {
	u, err := urlpkg.Parse(url)
	if err != nil {
		log.Error(err)
	}
	u.Host = removeEmptyPort(u.Host)
	r.URL = u
	r.Host = u.Host
	return r
}

// Given a string of the form "host", "host:port", or "[ipv6::address]:port",
// return true if the string includes a port.
func hasPort(s string) bool { return strings.LastIndex(s, ":") > strings.LastIndex(s, "]") }

// removeEmptyPort strips the empty port in ":port" to ""
// as mandated by RFC 3986 Section 6.2.3.
func removeEmptyPort(host string) string {
	if hasPort(host) {
		return strings.TrimSuffix(host, ":")
	}
	return host
}

func (r *replaceHttpRequest) SetMethod(method string) *replaceHttpRequest {
	r.Method = strings.ToUpper(method)
	return r
}

func (r *replaceHttpRequest) SetBody(body io.ReadCloser) *replaceHttpRequest {
	r.Body = body
	return r
}

func (r *replaceHttpRequest) SetContext(ctx context.Context) *replaceHttpRequest {
	stdr := (*http.Request)(r).WithContext(ctx)
	return (*replaceHttpRequest)(stdr)
}

func (r *replaceHttpRequest) StdHttpRequest() *http.Request {
	return (*http.Request)(r)
}
