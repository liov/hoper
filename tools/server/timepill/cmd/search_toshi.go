package main

import (
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	"github.com/actliboy/hoper/server/go/lib/utils/def/request"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"tools/timepill"
	"tools/toshi"
)

func main() {
	defer initialize.Start(&timepill.Conf, &timepill.Dao)()
	toshi.SetHost(timepill.Conf.Search.Host)
	//createIdx()
	//toshi.DisableLogger()

	req := &timepill.ListReq{
		ListReq: request.ListReq{
			PageReq: request.PageReq{PageNo: 1, PageSize: 100},
			SortReq: request.SortReq{SortField: "id", SortType: request.SortTypeASC},
		},
		RangeReq: request.RangeReq{
			RangeField: "id",
			RangeStart: toshimaxId(),
			RangeEnd:   nil,
			Include:    false,
		},
	}
	for {
		diaries, err := timepill.List(req)
		if err != nil {
			log.Error(err)
		}
		for _, diary := range diaries {
			toshi.AddDocument(timepill.DiaryIndex, true, diary.DiaryIndex())
		}
		if len(diaries) < req.PageSize {
			break
		}
		req.PageNo++
	}
}

func createIdx() {
	indexes := append(toshi.U64Indexes("id", "user_id", "notebook_id"), toshi.TextIndexes("notebook_subject", "content")...)
	indexes = append(indexes, toshi.DateIndex("created"))
	toshi.CreateIndex(timepill.DiaryIndex, indexes)
}

func exists(id int) bool {
	rep := toshi.Range[timepill.IndexDiary](timepill.DiaryIndex, map[string]*toshi.QueryRange{
		"id": {
			Gte: id,
			Lte: id,
		},
	}, 1)
	if rep.Hits > 0 {
		return true
	}
	return false
}

func fuzzy() {
	toshi.Term[timepill.IndexDiary](timepill.DiaryIndex, map[string]string{
		"created": "2010-03-18T13:03:48+00:00",
	}, 10)
	toshi.Fuzzy[timepill.IndexDiary](timepill.DiaryIndex, map[string]*toshi.QueryFuzzy{
		"content": {
			Value:         "今天",
			Distance:      0,
			Transposition: false,
		},
	}, 10)
}

func toshimaxId() int {
	rep := toshi.Range[timepill.IndexDiary](timepill.DiaryIndex, map[string]*toshi.QueryRange{
		"id": {},
	}, 1)
	if rep.Hits > 0 {
		return rep.Docs[0].Doc.Id
	}
	return 0
}
