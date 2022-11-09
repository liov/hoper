package main

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/actliboy/hoper/server/go/lib/initialize"
	"github.com/actliboy/hoper/server/go/lib/utils/conctrl"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"strconv"
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
	err := pro.Dao.Redis.SetNX(ctx, key, 510682, 0).Err()
	if err != nil {
		log.Error(err)
	}
	var id int
	err = pro.Dao.Redis.Get(ctx, key).Scan(&id)
	if err != nil {
		log.Error(err)
	}

	notFoundIds := make([]int, 0)
	timer := time.NewTicker(pro.Conf.Pro.Timer)
	for range timer.C {
		s, dir, _ := pro.Fetch(id)
		if s != nil {
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
			s.Each(func(i int, s *goquery.Selection) {
				if url, ok := s.Attr("file"); ok {
					conctrl.ReTry(5, func() error {
						return pro.Download(url, dir)
					})
				}
			})
			err := pro.Dao.Redis.Set(ctx, key, strconv.Itoa(id), 0).Err()
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
