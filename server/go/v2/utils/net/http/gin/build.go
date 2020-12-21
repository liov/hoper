package gin_build

import (
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/liov/hoper/go/v2/utils/configor"
	"github.com/liov/hoper/go/v2/utils/log"
)

func WithConfiguration(engine *gin.Engine, filename string) {
	url, err := url.Parse(filename)
	if err != nil {
		log.Fatalf("url解析错误：", err)
	}
	if url.Scheme == "nacos" {

	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Fatalf("配置文件不存在：", err)
	}

	err = configor.Load(engine, filename) //"./config/config.toml"
	if err != nil {
		log.Fatalf("配置错误: %v", err)
	}
}
