package route

import (
	"net/http"
	"testing"
)

func TestTrie(t *testing.T) {
	node := &node{}
	node.addRoute("/static/*filepath", &methodHandle{method: http.MethodGet})
	node.addRoute("/test/:id/:name", &methodHandle{method: http.MethodGet})
	node.addRoute("/", &methodHandle{method: http.MethodGet})
	node.addRoute("/apib", &methodHandle{method: http.MethodGet})
	node.addRoute("/api", &methodHandle{method: http.MethodGet})
	node.addRoute("/abc", &methodHandle{method: http.MethodGet})
	node.addRoute("/bcd", &methodHandle{method: http.MethodGet})
	node.addRoute("/test/:id", nil)
	node.addRoute("/test/:id", &methodHandle{method: http.MethodGet})
	node.addRoute("/test/:id", &methodHandle{method: http.MethodPost})
	node.addRoute("/abc/def", &methodHandle{method: http.MethodPost})
	node.addRoute("/test/:id/path/:id", &methodHandle{method: http.MethodGet})
	node.addRoute("/test/:id/path/:id", &methodHandle{method: http.MethodPost})
	node.addRoute("/test/:id/path/:id", &methodHandle{method: http.MethodPut})
	node.addRoute("/test/:id/path/path", &methodHandle{method: http.MethodGet})
	node.addRoute("/test/id/path/path/*path", nil)
	node.addRoute("/test/id/path/path/*path", &methodHandle{method: http.MethodPost})
	node.addRoute("/test/id/path/path/path", &methodHandle{method: http.MethodPost})
	//node.addRoute("/*filepath", &methodHandle{method: http.MethodGet})
	node.addRoute("/id", &methodHandle{method: http.MethodGet})
}
