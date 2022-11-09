package main

import (
	"github.com/actliboy/hoper/server/go/lib/initialize"
	"github.com/actliboy/hoper/server/go/lib/utils/generics/net/http/client/crawler"
	"tools/pro"
)

func main() {
	defer initialize.Start(&pro.Conf, &pro.Dao)()
	engine := crawler.NewEngine(pro.Conf.Pro.WorkCount)
	go normal(engine)
	engine.ErrHandler(pro.ErrorHandler())
	engine.Run()
}

func normal(engine *crawler.Engine) {
	start := 526399
	end := 532535
	for i := start; i <= end; i++ {
		req := pro.GetFetchReq(i)
		engine.BaseEngine.AddTask(engine.NewTask(req))

	}
}
