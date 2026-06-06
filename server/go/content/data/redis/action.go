package redis

import (
	"context"
	"strconv"

	redisi "github.com/hopeio/gox/database/redis"
	"github.com/hopeio/scaffold/errcode"
	"github.com/liov/hoper/server/go/protobuf/content"
)

func (d *ContentDao) HotCount(ctx context.Context, typ content.ContentType, refId uint64, changeCount float64) error {

	key := content.ContentType_name[int32(typ)][7:] + redisi.KeySortSet
	err := d.ZIncrBy(ctx, key, changeCount, strconv.FormatUint(refId, 10)).Err()
	if err != nil {
		return errcode.RedisErr.Wrap(err)
	}
	return nil
}

func (d *ContentDao) ActionCount(ctx context.Context, typ content.ContentType, action content.ActionType, refId uint64, changeCount float64) error {

	key := content.ContentType_name[int32(typ)][7:] + content.ActionType_name[int32(action)][6:] + redisi.KeySortSet
	err := d.ZIncrBy(ctx, key, changeCount, strconv.FormatUint(refId, 10)).Err()
	if err != nil {
		return errcode.RedisErr.Wrap(err)
	}
	return nil
}
