package main

import (
	"context"
	"github.com/actliboy/hoper/server/go/lib/initialize"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"time"
	"tools/timepill"
)

func main() {
	defer initialize.Start(&timepill.Conf, &timepill.Dao)()
	ctx := context.Background()
	key := "RecordByOrderUserID"
	err := timepill.Dao.Redis.SetNX(ctx, key, 1, 0).Err()
	if err != nil {
		log.Error(err)
	}
	var id int
	err = timepill.Dao.Redis.Get(ctx, key).Scan(&id)
	if err != nil {
		log.Error(err)
	}
	tc := time.NewTicker(time.Second * 2)
	for {
		if timepill.UserExists(id) {
			id++
			err = timepill.Dao.Redis.Incr(ctx, key).Err()
			if err != nil {
				log.Error(err)
			}
			continue
		}
		<-tc.C
		user := timepill.RecordUserById(id)
		if user.UserId != 0 {
			timepill.RecordUserDiaries(user)
		} else {
			time.Sleep(time.Second)
		}
		id++
		err = timepill.Dao.Redis.Incr(ctx, key).Err()
		if err != nil {
			log.Error(err)
		}
	}
}
