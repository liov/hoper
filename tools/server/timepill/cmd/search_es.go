package main

import (
	"context"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	v7 "github.com/actliboy/hoper/server/go/lib/tiga/initialize/elastic/v7"
	"github.com/actliboy/hoper/server/go/lib/utils/def/request"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"strconv"
	"tools/timepill"
)

var es v7.Es

func main() {
	defer initialize.Start(&timepill.Conf, &timepill.Dao)()
	es = timepill.Dao.Es
	ctx := context.Background()
	exists, err := es.IndexExists(timepill.DiaryIndex).Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		createIndex(ctx)
	}
	load(ctx)
}

func createIndex(ctx context.Context) {

	createIndex, err := es.CreateIndex(timepill.DiaryIndex).BodyString(timepill.Mapping).Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	if !createIndex.Acknowledged {
		// Not acknowledged
	}
}

func esmaxId(ctx context.Context) int {
	rep, _ := es.Search(timepill.DiaryIndex).Sort("id", false).Size(1).Do(ctx)
	if rep.TotalHits() > 0 {
		id, _ := strconv.Atoi(rep.Hits.Hits[0].Id)
		return id
	}
	return 0
}

func load(ctx context.Context) {
	req := &timepill.ListReq{
		ListReq: request.ListReq{
			PageReq: request.PageReq{PageNo: 1, PageSize: timepill.Conf.TimePill.PageSize},
			SortReq: request.SortReq{SortField: "id", SortType: request.SortTypeASC},
		},
		RangeReq: request.RangeReq{
			RangeField: "id",
			RangeStart: esmaxId(ctx),
			RangeEnd:   nil,
			Include:    false,
		},
	}
	index := es.Index().Index(timepill.DiaryIndex)
	for {
		req.PageSize = timepill.Conf.TimePill.PageSize
		if req.PageSize < 1 {
			req.PageSize = 10
		}
		diaries, err := timepill.Dao.ListDB(req)
		if err != nil {
			log.Error(err)
		}
		for i, diary := range diaries {
			_, err = index.Id(strconv.Itoa(diary.Id)).BodyJson(diary.DiaryIndex()).Do(ctx)
			if err != nil {
				log.Error(err)
			}
			if i == len(diaries)-1 {
				req.RangeStart = diary.Id
			}
		}
		if len(diaries) < req.PageSize {
			break
		}
	}
}
