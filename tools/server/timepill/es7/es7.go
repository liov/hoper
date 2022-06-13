package es7

import (
	"context"
	v7 "github.com/actliboy/hoper/server/go/lib/tiga/initialize/elastic/v7"
	"github.com/actliboy/hoper/server/go/lib/utils/def/request"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"strconv"
	"tools/timepill"
)

type Dao struct {
	ctx context.Context
	Es  v7.Es
}

func (dao *Dao) CreateIndexEs7() {
	exists, err := dao.Es.IndexExists(timepill.DiaryIndex).Do(dao.ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		createIndex, err := dao.Es.CreateIndex(timepill.DiaryIndex).BodyString(timepill.Mapping).Do(dao.ctx)
		if err != nil {
			// Handle error
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}
}

func (dao *Dao) MaxIdEs7() int {
	rep, _ := dao.Es.Search(timepill.DiaryIndex).Sort("id", false).Size(1).Do(dao.ctx)
	if rep.TotalHits() > 0 {
		id, _ := strconv.Atoi(rep.Hits.Hits[0].Id)
		return id
	}
	return 0
}

func (dao *Dao) LoadES7(pigeSize int) {
	req := &timepill.ListReq{
		ListReq: request.ListReq{
			PageReq: request.PageReq{PageNo: 1, PageSize: pigeSize},
			SortReq: request.SortReq{SortField: "id", SortType: request.SortTypeASC},
		},
		RangeReq: request.RangeReq{
			RangeField: "id",
			RangeStart: dao.MaxIdEs7(),
			RangeEnd:   nil,
			Include:    false,
		},
	}
	index := dao.Es.Index().Index(timepill.DiaryIndex)
	for {
		req.PageSize = pigeSize
		if req.PageSize < 1 {
			req.PageSize = 10
		}
		diaries, err := timepill.Dao.DBDao(dao.ctx).ListDB(req)
		if err != nil {
			log.Error(err)
		}
		for i, diary := range diaries {
			_, err = index.Id(strconv.Itoa(diary.Id)).BodyJson(diary.DiaryIndex()).Do(dao.ctx)
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
