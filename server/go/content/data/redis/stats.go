package redis

import (
	"strconv"

	"github.com/hopeio/scaffold/errcode"
	"github.com/liov/hoper/server/go/content/model"
)

func (d *ContentDao) UserContentEdit(userId uint64, field string, value interface{}) error {
	ctx := d.Context()
	key := model.UserContentCountKey + strconv.FormatUint(userId, 10)

	err := d.HSet(ctx, key, field, value).Err()
	if err != nil {
		return errcode.RedisErr.Wrap(err)
	}
	return nil
}
