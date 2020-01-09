package iris_build

import (
	"os"

	"github.com/jinzhu/configor"
	"github.com/kataras/iris/v12"
	"github.com/liov/hoper/go/v2/utils/log"
)

type Config struct {
	Iris iris.Configuration
}

func Build(app *iris.Application, filename string) {
	//将来这块用配置中心
	if filename != "" {
		WithConfiguration(app, filename)
	}
	if err := app.Build(); err != nil {
		log.Fatal(err)
	}
}

func Configuration(filename string) *iris.Configuration {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Fatalf("配置文件不存在：", err)
	}
	var c Config
	err := configor.New(&configor.Config{Debug: false}).
		Load(&c, filename) //"./config/config.toml"
	if err != nil {
		log.Fatalf("配置错误: %v", err)
	}
	return &c.Iris
}

func WithConfiguration(mux *iris.Application, filename string) {
	mux.Configure(iris.WithConfiguration(*Configuration(filename)))
}
