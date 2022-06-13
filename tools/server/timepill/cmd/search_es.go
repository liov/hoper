package main

import (
	"context"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	"tools/timepill"
)

func main() {
	defer initialize.Start(&timepill.Conf, &timepill.Dao)()
	timepill.Dao.EsDao(context.Background()).LoadEs8()
}
