package redis

import (
	"github.com/hopeio/tiga/protobuf/errorcode"
	"github.com/liov/hoper/server/go/content/model"
)

func (d *ContentRedisDao) UserContentEdit(field string, value interface{}) error {
	ctxi := d
	ctx := ctxi.Context.Context()
	key := model.UserContentCountKey + ctxi.AuthID

	err := d.conn.HSet(ctx, key, field, value).Err()
	if err != nil {
		return ctxi.ErrorLog(errorcode.RedisErr, err, "RedisUserInfoEdit")
	}
	return nil
}
