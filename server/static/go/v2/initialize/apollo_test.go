package initialize

import (
	"testing"

	"github.com/liov/hoper/go/v2/utils/log"
)

func TestApolloConfig_Generate(t *testing.T) {
	type Config struct {
		A int    `json:"a"`
		B string `json:"b"`
	}
	type Foo struct {
		Bar Config
	}
	var foo Foo
	var conf = map[string]string{"baR": "A = 1\nB = \"哈哈哈\""}
	ApolloConfigEnable(&foo, conf)
	log.Info(foo)
}
