package main

import (
	"context"
	"fmt"
	"github.com/liov/hoper/server/go/lib/initialize"
	"github.com/liov/hoper/server/go/lib/utils/log"
	"time"
	"tools/clawer/timepill"
)

func main() {
	defer initialize.Start(&timepill.Conf, &timepill.Dao)()
	fmt.Println(timepill.Token)
	ctx := context.Background()
	key := "RecordByOrderNoteBookID"
	err := timepill.Dao.Redis.SetNX(ctx, key, 1, 0).Err()
	if err != nil {
		log.Error(err)
	}
	var id int
	err = timepill.Dao.Redis.Get(ctx, key).Scan(&id)
	if err != nil {
		log.Error(err)
	}
	var maxId int
	err = timepill.Dao.Hoper.Raw(`SELECT MAX(id) FROM "note_book"`).Row().Scan(&maxId)
	if err != nil {
		log.Error(err)
	}
	var continuouZeroId int
	tc := time.NewTicker(time.Second)
	for {
		var exists bool
		err = timepill.Dao.Hoper.Raw(`SELECT EXISTS(SELECT id FROM "note_book" WHERE id = ? AND expired < '2022-01-01' LIMIT 1)`, id).Row().Scan(&exists)
		if err != nil {
			log.Error(err)
		}

		if exists {
			id++
			err = timepill.Dao.Redis.Incr(ctx, key).Err()
			if err != nil {
				log.Error(err)
			}
			continue
		}
		<-tc.C
		notebook := timepill.RecordByNoteBookId(id)
		if notebook.Id == 0 {
			continuouZeroId++
			if continuouZeroId == 100 && id > maxId {
				tc.Stop()
				return
			}
		} else {
			continuouZeroId = 0
		}
		id++
		err = timepill.Dao.Redis.Incr(ctx, key).Err()
		if err != nil {
			log.Error(err)
		}
	}
	//timepill.RecordByNoteBookId(873)
}
