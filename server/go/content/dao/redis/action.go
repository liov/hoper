package redis

import (
	"github.com/actliboy/hoper/server/go/protobuf/content"
	"github.com/hopeio/lemon/protobuf/errorcode"
	redisi "github.com/hopeio/lemon/utils/dao/redis"
	"strconv"
)

func (d *ContentRedisDao) HotCount(typ content.ContentType, refId uint64, changeCount float64) error {
	ctxi := d.Context
	key := content.ContentType_name[int32(typ)][7:] + redisi.SortSet
	err := d.conn.ZIncrBy(ctxi.Context, key, changeCount, strconv.FormatUint(refId, 10)).Err()
	if err != nil {
		return ctxi.ErrorLog(errorcode.RedisErr, err, "HotCountRedis")
	}
	return nil
}

func (d *ContentRedisDao) ActionCount(typ content.ContentType, action content.ActionType, refId uint64, changeCount float64) error {
	ctxi := d.Context
	key := content.ContentType_name[int32(typ)][7:] + content.ActionType_name[int32(action)][6:] + redisi.SortSet
	err := d.conn.ZIncrBy(ctxi.Context, key, changeCount, strconv.FormatUint(refId, 10)).Err()
	if err != nil {
		return ctxi.ErrorLog(errorcode.RedisErr, err, "ActionCountRedis")
	}
	return nil
}
