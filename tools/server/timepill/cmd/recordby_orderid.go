package main

import (
	"context"
	"github.com/actliboy/hoper/server/go/lib/initialize"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"tools/timepill"
)

func main() {
	defer initialize.Start(&timepill.Conf, &timepill.Dao)()

	ctx := context.Background()

	maxid, _ := timepill.Dao.DBDao(ctx).MaxDiaryId()

	key := "RecordByOrderID2"
	err := timepill.Dao.Redis.SetNX(ctx, key, 1, 0).Err()
	if err != nil {
		log.Error(err)
	}
	var id int
	err = timepill.Dao.Redis.Get(ctx, key).Scan(&id)
	if err != nil {
		log.Error(err)
	}
	for {
		if timepill.DiaryExists(id) {
			id++
			err = timepill.Dao.Redis.Incr(ctx, key).Err()
			if err != nil {
				log.Error(err)
			}
			timepill.RecordCommentWithJudge(id)
			continue
		}
		timepill.RecordDiaryById(id)
		id++
		err = timepill.Dao.Redis.Incr(ctx, key).Err()
		if err != nil {
			log.Error(err)
		}
		if id >= maxid {
			break
		}
	}
}
