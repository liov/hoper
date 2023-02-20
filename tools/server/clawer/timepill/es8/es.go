package es8

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/liov/hoper/server/go/lib/utils/def/request"
	"github.com/liov/hoper/server/go/lib/utils/io/reader"
	"github.com/liov/hoper/server/go/lib/utils/log"
	_type "github.com/liov/hoper/server/go/lib_v2/utils/dao/db/type"
	utilv8 "github.com/liov/hoper/server/go/lib_v2/utils/dao/elastic/v8"
	"net/http"
	"strconv"
	"strings"
	"tools/clawer/timepill"
	"tools/clawer/timepill/model"
)

const DiaryIndex = "diary"

const Mapping = `
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

type EsDao struct {
	ctx context.Context
	Es8 *elasticsearch.Client
}

func NewEsDao(ctx context.Context, es8 *elasticsearch.Client) *EsDao {
	return &EsDao{ctx: ctx, Es8: es8}
}

func (dao *EsDao) MaxIdEs8() int {
	size := 1
	req := esapi.SearchRequest{
		Index: []string{DiaryIndex},
		Sort:  []string{"id:desc"},
		Size:  &size,
	}

	resp, err := utilv8.GetSearchResponse[model.IndexDiary](req.Do(dao.ctx, dao.Es8))
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

func (dao *EsDao) LoadEs8() {
	req := &_type.ListReq[int]{
		PageSortReq: request.PageSortReq{
			PageReq: request.PageReq{PageNo: 1, PageSize: timepill.Conf.TimePill.PageSize},
			SortReq: &request.SortReq{SortField: "id", SortType: request.SortTypeASC},
		},
		RangeReq: &request.RangeReq[int]{
			RangeField: "id",
			RangeStart: dao.MaxIdEs8(),
			RangeEnd:   0,
			Include:    false,
		},
	}

	for {
		req.PageSize = timepill.Conf.TimePill.PageSize
		if req.PageSize < 1 {
			req.PageSize = 10
		}
		diaries, err := timepill.Dao.DBDao(dao.ctx).List(req)
		if err != nil {
			log.Error(err)
		}
		for i, diary := range diaries {
			body, _ := json.Marshal(diary.IndexDiary())
			esreq := esapi.CreateRequest{
				Index:      DiaryIndex,
				DocumentID: strconv.Itoa(diary.Id),
				Body:       bytes.NewReader(body),
			}
			resp, err := esreq.Do(dao.ctx, dao.Es8)
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

func (dao *EsDao) CreateIndexEs8() {
	resp, err := dao.Es8.Indices.Exists([]string{DiaryIndex})
	if err != nil {
		// Handle error
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		log.Info("index not found")
		_, err := dao.Es8.Indices.Create(DiaryIndex, dao.Es8.Indices.Create.WithBody(strings.NewReader(Mapping)))
		if err != nil {
			panic(err)
		}
		resp.Body.Close()
		log.Info("index created")
	} else {
		log.Info("index found")
	}
}
