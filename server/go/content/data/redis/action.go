package redis

import (
	"github.com/hopeio/scaffold/errcode"
	redisi "github.com/hopeio/utils/datax/redis"
	"github.com/liov/hoper/server/go/protobuf/content"
	"strconv"
)

func (d *ContentDao) HotCount(typ content.ContentType, refId uint64, changeCount float64) error {
	ctxi := d.Context
	key := content.ContentType_name[int32(typ)][7:] + redisi.KeySortSet
	err := d.conn.ZIncrBy(ctxi.Base(), key, changeCount, strconv.FormatUint(refId, 10)).Err()
	if err != nil {
		return ctxi.RespErrorLog(errcode.RedisErr, err, "HotCountRedis")
	}
	return nil
}

func (d *ContentDao) ActionCount(typ content.ContentType, action content.ActionType, refId uint64, changeCount float64) error {
	ctxi := d.Context
	key := content.ContentType_name[int32(typ)][7:] + content.ActionType_name[int32(action)][6:] + redisi.KeySortSet
	err := d.conn.ZIncrBy(ctxi.Base(), key, changeCount, strconv.FormatUint(refId, 10)).Err()
	if err != nil {
		return ctxi.RespErrorLog(errcode.RedisErr, err, "ActionCountRedis")
	}
	return nil
}
