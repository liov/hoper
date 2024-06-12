package redis

import (
	"github.com/hopeio/cherry/protobuf/errcode"
	"github.com/liov/hoper/server/go/content/model"
)

func (d *ContentDao) UserContentEdit(field string, value interface{}) error {
	ctxi := d
	ctx := ctxi.Context.Context()
	key := model.UserContentCountKey + ctxi.AuthID

	err := d.conn.HSet(ctx, key, field, value).Err()
	if err != nil {
		return ctxi.ErrorLog(errcode.RedisErr, err, "RedisUserInfoEdit")
	}
	return nil
}
