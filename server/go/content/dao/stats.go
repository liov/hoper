package dao

import (
	"github.com/liov/hoper/v2/content/model"
	"github.com/liov/hoper/v2/protobuf/utils/errorcode"
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
