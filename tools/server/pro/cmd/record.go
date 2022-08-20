package main

import (
	"context"
	"github.com/actliboy/hoper/server/go/lib/initialize"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"strings"
	"time"
	"tools/pro"
)

func main() {
	defer initialize.Start(&pro.Conf, &pro.Dao)()

	record()
}

func record() {
	ctx := context.Background()

	key := "Pro_RecordByOrderID"
	err := pro.Dao.Redis.SetNX(ctx, key, 492110, 0).Err()
	if err != nil {
		log.Error(err)
	}
	var id int
	err = pro.Dao.Redis.Get(ctx, key).Scan(&id)
	if err != nil {
		log.Error(err)
	}
	sd := pro.NewSpeed(pro.Conf.Pro.Loop)
	sd.FailPic.Record("fail_pic_")
	sd.FailDB.Record("fail_db_")
	notFoundIds := make([]int, 0)
	timer := time.NewTicker(pro.Conf.Pro.Timer)
	for range timer.C {
		tid := pro.Record(id, sd)
		if tid != "" {
			if len(notFoundIds) > 0 {
				for _, id := range notFoundIds {
					invalidPost := &pro.Post{TId: id, Status: 2}
					err := pro.Dao.DB.Save(invalidPost).Error
					if err != nil && !strings.HasPrefix(err.Error(), "ERROR: duplicate key") {
						log.Error(err)
					}
				}
				notFoundIds = notFoundIds[:0]
			}
			err := pro.Dao.Redis.Set(ctx, key, tid, 0).Err()
			if err != nil {
				log.Error(err)
			}
		} else {
			notFoundIds = append(notFoundIds, id)
			if len(notFoundIds) >= 20 {
				err = pro.Dao.Redis.Get(ctx, key).Scan(&id)
				if err != nil {
					log.Error(err)
				}
				notFoundIds = notFoundIds[:0]
			}
		}
		id++
	}
}
