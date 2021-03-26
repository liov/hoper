package _pick

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/json-iterator/go/extra"
	"github.com/liov/hoper/go/v2/utils/encoding/json"
)

func TestTrie(t *testing.T) {
	node := &node{}
	node.addRoute("/static/*filepath", &methodHandle{method: http.MethodGet,httpHandler: http.NotFoundHandler()})
	//node.addRoute("/test/:id/:name", &methodHandle{method: http.MethodGet,httpHandler: http.NotFoundHandler()})
	node.addRoute("/", &methodHandle{method: http.MethodGet,httpHandler: http.NotFoundHandler()})
	node.addRoute("/apib", &methodHandle{method: http.MethodGet,httpHandler: http.NotFoundHandler()})
	node.addRoute("/api", &methodHandle{method: http.MethodGet,httpHandler: http.NotFoundHandler()})
	node.addRoute("/abc", &methodHandle{method: http.MethodGet,httpHandler: http.NotFoundHandler()})
	node.addRoute("/bcd", &methodHandle{method: http.MethodGet,httpHandler: http.NotFoundHandler()})
	//node.addRoute("/test/:id", &methodHandle{method: http.MethodGet,httpHandler: http.NotFoundHandler()})
	//node.addRoute("/test/:id", &methodHandle{method: http.MethodPost,httpHandler: http.NotFoundHandler()})
	node.addRoute("/abc/def", &methodHandle{method: http.MethodPost,httpHandler: http.NotFoundHandler()})
	/*node.addRoute("/test/:id/path/:id", &methodHandle{method: http.MethodGet,httpHandler: http.NotFoundHandler()})
	node.addRoute("/test/:id/path/:id", &methodHandle{method: http.MethodPost,httpHandler: http.NotFoundHandler()})
	node.addRoute("/test/:id/path/:id", &methodHandle{method: http.MethodPut,httpHandler: http.NotFoundHandler()})
	node.addRoute("/test/:id/path/path", &methodHandle{method: http.MethodGet,httpHandler: http.NotFoundHandler()})*/
	node.addRoute("/test/id/path/path/*path", &methodHandle{method: http.MethodPost,httpHandler: http.NotFoundHandler()})
	//node.addRoute("/test/id/path/path/path", &methodHandle{method: http.MethodPost,httpHandler: http.NotFoundHandler()})
//	node.addRoute("/*filepath", &methodHandle{method: http.MethodGet,httpHandler: http.NotFoundHandler()})
	node.addRoute("/id", &methodHandle{method: http.MethodGet,httpHandler: http.NotFoundHandler()})
	extra.SupportPrivateFields()
	data,err:=json.Marshal(node)
	if err != nil {
		t.Log(err)
	}
	fmt.Println(string(data))
	fmt.Printf("%#v\n", node)
}
