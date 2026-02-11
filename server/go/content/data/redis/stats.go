package redis

import (
	"context"
	"strconv"

	"github.com/hopeio/scaffold/errcode"
	"github.com/liov/hoper/server/go/content/model"
)

func (d *ContentDao) UserContentEdit(ctx context.Context, userId uint64, field string, value interface{}) error {

	key := model.UserContentCountKey + strconv.FormatUint(userId, 10)

	err := d.HSet(ctx, key, field, value).Err()
	if err != nil {
		return errcode.RedisErr.Wrap(err)
	}
	return nil
}
