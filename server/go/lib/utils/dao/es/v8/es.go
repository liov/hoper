package v8

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/actliboy/hoper/server/go/lib/utils/io/reader"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"net/http"
)

type SearchResponse[T any] struct {
	Took     int     `json:"took"`
	TimedOut bool    `json:"timed_out"`
	Shards   Shards  `json:"_shards"`
	Hits     Hits[T] `json:"hits"`
}

type Shards struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}

type Hits[T any] struct {
	Total    Total    `json:"total"`
	MaxScore any      `json:"max_score"`
	Hits     []Hit[T] `json:"hits"`
}

type Hit[T any] struct {
	Index  string      `json:"_index"`
	Id     string      `json:"_id"`
	Score  interface{} `json:"_score"`
	Source T           `json:"_source"`
	Sort   []int       `json:"sort"`
}

type Total struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}

func GetResponse[T any](response *esapi.Response, err error) (*T, error) {
	if err != nil {
		return nil, err
	}
	bytes, err := reader.ReadCloser(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, errors.New(string(bytes))
	}
	var res T
	err = json.Unmarshal(bytes, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func GetSearchResponse[T any](response *esapi.Response, err error) (*SearchResponse[T], error) {
	return GetResponse[SearchResponse[T]](response, err)
}

func CreateDocument[T any](ctx context.Context, es *elasticsearch.Client, index, id string, obj T) error {
	body, _ := json.Marshal(obj)
	esreq := esapi.CreateRequest{
		Index:      index,
		DocumentID: id,
		Body:       bytes.NewReader(body),
	}
	resp, err := esreq.Do(ctx, es)
	if err != nil {
		return err
	}
	return resp.Body.Close()
}
