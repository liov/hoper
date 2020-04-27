package trie

import (
	"net/http"
	"testing"
)

func TestTrie(t *testing.T) {
	node := &node{}
	node.addRoute("/*filepath", &methodHandle{method: http.MethodGet})
	node.addRoute("/:id/:name", &methodHandle{method: http.MethodGet})
	node.addRoute("/", &methodHandle{method: http.MethodGet})
	node.addRoute("/apib", &methodHandle{method: http.MethodGet})
	node.addRoute("/api", &methodHandle{method: http.MethodGet})
	node.addRoute("/abc", &methodHandle{method: http.MethodGet})
	node.addRoute("/bcd", &methodHandle{method: http.MethodGet})
	node.addRoute("/:id", nil)
	node.addRoute("/:id", &methodHandle{method: http.MethodGet})
	node.addRoute("/:id", &methodHandle{method: http.MethodPost})
	node.addRoute("/abc/def", &methodHandle{method: http.MethodPost})
	node.addRoute("/:id/path/:id", &methodHandle{method: http.MethodGet})
	node.addRoute("/:id/path/:id", &methodHandle{method: http.MethodPost})
	node.addRoute("/:id/path/:id", &methodHandle{method: http.MethodPut})
	node.addRoute("/:id/path/path", &methodHandle{method: http.MethodGet})
	node.addRoute("/id/path/path/*path", &methodHandle{method: http.MethodGet})
	//node.addRoute("/*filepath", &methodHandle{method: http.MethodGet})
	node.addRoute("/id", &methodHandle{method: http.MethodGet})
}
