package main

import (
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	"github.com/actliboy/hoper/server/go/lib/tiga/pick"
	"tools/timepill"
)

func main() {
	defer initialize.Start(&timepill.Conf, &timepill.Dao, "Kafka")()
	pick.RegisterService(userservice.GetUserService(), contentervice.GetMomentService())
}
