package main

import (
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/utils/dao/db/get"
)

var userMod = []interface{}{
	&model.User{},
	&model.UserExtend{},
	&model.UserActionLog{},
	&model.UserBannedLog{},
	&model.UserFollow{},
	&model.UserScoreLog{},
	&model.UserFollowLog{},
	&model.Resume{},
}

func main() {
	get.OrmDB.DropTable(userMod...)
	get.OrmDB.CreateTable(userMod...)
}
