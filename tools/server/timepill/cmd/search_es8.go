package main

import (
	"bytes"
	"context"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	v8 "github.com/actliboy/hoper/server/go/lib/tiga/initialize/elastic/v8"
	utilv8 "github.com/actliboy/hoper/server/go/lib/utils/dao/es/v8"
	"github.com/actliboy/hoper/server/go/lib/utils/def/request"
	"github.com/actliboy/hoper/server/go/lib/utils/encoding/json"
	"github.com/actliboy/hoper/server/go/lib/utils/io/reader"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"net/http"
	"strconv"
	"strings"
	"tools/timepill"
)

var es8 v8.Es

func main() {
	defer initialize.Start(&timepill.Conf, &timepill.Dao)()
	es8 = timepill.Dao.Es8
	resp, err := es8.Indices.Exists([]string{timepill.DiaryIndex})
	if err != nil {
		// Handle error
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		log.Info("index not found")
		_, err := es8.Indices.Create(timepill.DiaryIndex, es8.Indices.Create.WithBody(strings.NewReader(timepill.Mapping)))
		if err != nil {
			panic(err)
		}
		bytes, err := reader.ReadCloser(resp.Body)
		log.Info(string(bytes))
		log.Info("index created")
	} else {
		log.Info("index found")
	}
	load(context.Background())
}

func esmaxId(ctx context.Context) int {
	size := 1
	req := esapi.SearchRequest{
		Index: []string{timepill.DiaryIndex},
		Sort:  []string{"id:desc"},
		Size:  &size,
	}

	resp, err := utilv8.GetSearchResponse[timepill.IndexDiary](req.Do(ctx, es8))
	if err != nil {
		log.Error(err)
		return 0
	}
	if resp.Hits.Total.Value > 0 {
		id, _ := strconv.Atoi(resp.Hits.Hits[0].Id)
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
			body, _ := json.Marshal(diary.DiaryIndex())
			esreq := esapi.CreateRequest{
				Index:      timepill.DiaryIndex,
				DocumentID: strconv.Itoa(diary.Id),
				Body:       bytes.NewReader(body),
			}
			resp, err := esreq.Do(ctx, es8)
			if err != nil {
				// Handle error
				log.Error(err)
			}
			bytes, err := reader.ReadCloser(resp.Body)
			log.Info(string(bytes))
			if i == len(diaries)-1 {
				req.RangeStart = diary.Id
			}
		}
		if len(diaries) < req.PageSize {
			break
		}
	}
}
