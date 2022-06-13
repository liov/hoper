package redis

import (
	"github.com/actliboy/hoper/server/go/lib/protobuf/errorcode"
	"github.com/actliboy/hoper/server/go/mod/content/model"
)

func (d *ContentRedisDao) UserContentEdit(field string, value interface{}) error {
	ctxi := d
	ctx := ctxi.Context
	key := model.UserContentCountKey + ctxi.IdStr

	err := d.conn.HSet(ctx, key, field, value).Err()
	if err != nil {
		return ctxi.ErrorLog(errorcode.RedisErr, err, "RedisUserInfoEdit")
	}
	return nil
}
