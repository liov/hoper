package main

import (
	"github.com/liov/hoper/server/go/lib/initialize"
	"github.com/liov/hoper/server/go/lib/utils/net/http/client/crawler"
	"tools/clawer/pro"
)

func main() {
	defer initialize.Start(&pro.Conf, &pro.Dao)()
	engine := crawler.NewEngine(pro.Conf.Pro.WorkCount)
	go normal(engine)
	engine.ErrHandler(pro.ErrorHandler())
	engine.Run()
}

func normal(engine *crawler.Engine) {
	start := 552501
	end := 552600
	for i := start; i <= end; i++ {
		req := pro.GetFetchReq(i)
		engine.BaseEngine.AddTask(engine.BaseTask(req))

	}
}
