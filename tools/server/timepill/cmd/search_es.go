package main

import (
	"context"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	"github.com/actliboy/hoper/server/go/lib/utils/def/request"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"strconv"
	"tools/timepill"
)

func main() {
	defer initialize.Start(&timepill.Conf, &timepill.Dao)()
	ctx := context.Background()
	exists, err := timepill.Dao.Es.IndexExists(timepill.DiaryIndex).Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		createIndex(ctx)
	}
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
	index := timepill.Dao.Es.Index().Index(timepill.DiaryIndex)
	for {
		req.PageSize = timepill.Conf.TimePill.PageSize
		if req.PageSize < 1 {
			req.PageSize = 10
		}
		diaries, err := timepill.List(req)
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

func createIndex(ctx context.Context) {
	const mapping = `
{
    "settings": {
        "number_of_shards": 1,
        "number_of_replicas": 0
    },
    "mappings": {
        "properties": {
			"id": {
                "type": "long"
            },
            "user_id": {
                "type": "keyword",
                "fields": {
                    "raw": {
                        "type": "long"
                    }
                }
            },
            "notebook_id": {
                "type": "keyword",
                "fields": {
                    "raw": {
                        "type": "long"
                    }
                }
            },
            "notebook_subject": {
                "type": "text",
                "analyzer": "ik_max_word",
                "search_analyzer": "ik_smart",
                "store": true,
                "fielddata": true
            },
            "content": {
                "type": "text",
                "analyzer": "ik_max_word",
                "search_analyzer": "ik_smart"
            },
            "created": {
                "type": "date"
            }
        }
    }
}`
	createIndex, err := timepill.Dao.Es.CreateIndex(timepill.DiaryIndex).BodyString(mapping).Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	if !createIndex.Acknowledged {
		// Not acknowledged
	}
}

func esmaxId(ctx context.Context) int {
	rep, _ := timepill.Dao.Es.Search(timepill.DiaryIndex).Sort("id", false).Size(1).Do(ctx)
	if rep.TotalHits() > 0 {
		id, _ := strconv.Atoi(rep.Hits.Hits[0].Id)
		return id
	}
	return 0
}
