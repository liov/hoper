package global

import (
	"github.com/hopeio/gox/idgen"
	"github.com/liov/hoper/server/go/global"
)

var Dao = global.Global.Dao
var Conf = global.Global.Config

var SF *idgen.Snowflake
