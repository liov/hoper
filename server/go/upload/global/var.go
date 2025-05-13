package global

import (
	"github.com/liov/hoper/server/go/global"
)

var Dao = global.Global.Dao
var Conf = &config{}

func init() {
	global.Global.Unmarshal(Conf)
}
