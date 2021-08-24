package dao

import (
	"github.com/liov/hoper/server/go/lib/protobuf/errorcode"
	"github.com/liov/hoper/server/go/mod/content/model"
)

func (d *contentDao) UserContentEditRedis(field string, value interface{}) error {
	ctxi := d
	ctx := ctxi.Context
	key := model.UserContentCountKey + ctxi.IdStr

	err := Dao.Redis.HSet(ctx, key, field, value).Err()
	if err != nil {
		return ctxi.ErrorLog(errorcode.RedisErr, err, "RedisUserInfoEdit")
	}
	return nil
}
