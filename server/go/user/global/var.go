package global

import "github.com/hopeio/initialize"

var Global = initialize.NewGlobal[*config, *dao]()
var Dao *dao = Global.Dao
var Conf *config = Global.Config
