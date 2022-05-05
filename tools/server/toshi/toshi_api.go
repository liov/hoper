package toshi

import (
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client"
	"golang.org/x/exp/constraints"
	"net/http"
)

var toshiRequest = client.NewEasyRequest() //.SetLogger(nil)

var host = "http://es.liov.xyz/"

type Config struct {
	Host string
}

func SetHost(newHost string) {
	host = newHost
}

func DisableLogger() {
	toshiRequest.SetLogger(nil)
}

type ResponseBody struct {
	Status  int         `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func toshiCall(api, method string, param, response any) {
	err := toshiRequest.CompleteDo(host+api, method, param, response)
	if err != nil {
		log.Error(err)
	}
}

func get(api string, param, response any) {
	toshiCall(api, http.MethodGet, param, response)
}
func put(api string, param, response any) {
	toshiCall(api, http.MethodPut, param, response)
}
func post(api string, param, response any) {
	toshiCall(api, http.MethodPost, param, response)
}
func delete(api string, param, response any) {
	toshiCall(api, http.MethodDelete, param, response)
}

type CreateIndexReq struct {
	Name    string                `json:"name"`
	Type    string                `json:"type"`
	Options CreateIndexReqOptions `json:"options"`
}

type CreateIndexReqOptions struct {
	Indexing *CreateIndexReqOptionsIndex `json:"indexing,omitempty"`
	Fast     string                      `json:"fast,omitempty"`
	Stored   bool                        `json:"stored"`
	Indexed  bool                        `json:"indexed,omitempty"`
}
type CreateIndexReqOptionsIndex struct {
	Record    string `json:"record"`
	Tokenizer string `json:"tokenizer"`
}

func CreateIndex(index string, indexes []*CreateIndexReq) {
	put(index+"/_create", indexes, nil)
}

type Id interface {
	~uint | ~int | ~int32 | ~uint32 | ~uint64 | ~int64
}

func U64Indexes(indexes ...string) []*CreateIndexReq {
	var indexesReq []*CreateIndexReq
	for i := range indexes {
		indexesReq = append(indexesReq, U64Indexe(indexes[i]))
	}
	return indexesReq
}

func U64Indexe(index string) *CreateIndexReq {
	indexReq := &CreateIndexReq{
		Name: index,
		Type: "u64",
		Options: CreateIndexReqOptions{
			Stored:  true,
			Indexed: true,
		},
	}
	if index == "id" {
		indexReq.Options.Fast = "single"
	}
	return indexReq
}

func TextIndexes(indexes ...string) []*CreateIndexReq {
	var indexesReq []*CreateIndexReq
	for i := range indexes {
		indexesReq = append(indexesReq, &CreateIndexReq{
			Name: indexes[i],
			Type: "text",
			Options: CreateIndexReqOptions{
				Stored: true,
				Indexing: &CreateIndexReqOptionsIndex{
					Record:    "position",
					Tokenizer: "CANG_JIE",
				},
			},
		})
	}
	return indexesReq
}

func DateIndexes(indexes ...string) []*CreateIndexReq {
	var indexesReq []*CreateIndexReq
	for i := range indexes {
		indexesReq = append(indexesReq, DateIndex(indexes[i]))
	}
	return indexesReq
}

func DateIndex(index string) *CreateIndexReq {
	return &CreateIndexReq{
		Name: index,
		Type: "date",
		Options: CreateIndexReqOptions{
			Stored:  true,
			Indexed: true,
		},
	}
}

type AddDocumentReq struct {
	Options  Options `json:"options"`
	Document any     `json:"document"`
}
type Options struct {
	Commit bool `json:"commit"`
}

func AddDocument(index string, commit bool, document any) {
	put(index+"/", &AddDocumentReq{
		Options:  Options{Commit: commit},
		Document: document,
	}, nil)
}

type QueryReq[T Word] struct {
	Query  QueryQuery[T] `json:"query"`
	Limit  int           `json:"limit"`
	SortBy string        `json:"sort_by,omitempty"`
}

type Word interface {
	constraints.Integer | constraints.Float | ~string | ~struct{}
}

type QueryQuery[T Word] struct {
	Term   map[string]string       `json:"term,omitempty"`
	Fuzzy  map[string]*QueryFuzzy  `json:"fuzzy,omitempty"`
	Phrase map[string]*QueryPhrase `json:"phrase,omitempty"`
	Range  map[string]*QueryRange  `json:"range,omitempty"`
	Bool   *QueryBool[T]           `json:"bool,omitempty"`
}

type QueryFuzzy struct {
	Value         string `json:"value"`
	Distance      int    `json:"distance"`
	Transposition bool   `json:"transposition"`
}

type QueryPhrase struct {
	Terms   []string `json:"terms"`
	Offsets []int    `json:"offsets,omitempty"`
}

type QueryBool[T Word] struct {
	Must               []QueryQuery[T] `json:"must"`
	MustNot            []QueryQuery[T] `json:"must_not"`
	Should             []QueryQuery[T] `json:"should"`
	MinimumShouldMatch int             `json:"minimum_should_match"`
}

type QueryRange struct {
	Gte int `json:"gte,omitempty"`
	Lte int `json:"lte,omitempty"`
	Gt  int `json:"gt,omitempty"`
	Lt  int `json:"lt,omitempty"`
}

type QueryRep[T any] struct {
	Hits   int           `json:"hits"`
	Docs   []*Doc[T]     `json:"docs"`
	Facets []interface{} `json:"facets"`
}

type Doc[T any] struct {
	Score float64 `json:"score"`
	Doc   T       `json:"doc"`
}

func Query[T Word, V any](index string, query *QueryReq[T]) *QueryRep[V] {
	var resp QueryRep[V]
	post(index+"/", query, &resp)
	return &resp
}

func Term[V any](index string, query map[string]string, limit int) *QueryRep[V] {
	var resp QueryRep[V]
	post(index+"/", &QueryReq[struct{}]{
		Query: QueryQuery[struct{}]{
			Term: query,
		},
		Limit:  limit,
		SortBy: "id",
	}, &resp)
	return &resp
}

func Fuzzy[V any](index string, query map[string]*QueryFuzzy, limit int) *QueryRep[V] {
	var resp QueryRep[V]
	post(index+"/", &QueryReq[struct{}]{
		Query: QueryQuery[struct{}]{
			Fuzzy: query,
		},
		Limit:  limit,
		SortBy: "id",
	}, &resp)
	return &resp
}

func Range[V any](index string, query map[string]*QueryRange, limit int) *QueryRep[V] {
	var resp QueryRep[V]
	post(index+"/", &QueryReq[struct{}]{
		Query: QueryQuery[struct{}]{
			Range: query,
		},
		Limit:  limit,
		SortBy: "id",
	}, &resp)
	return &resp
}

type DeleteDocReq struct {
	Options Options           `json:"options"`
	Term    map[string]string `json:"term,omitempty"`
}

type DocsAffected struct {
	DocsAffected int `json:"docs_affected"`
}

func Delete(index string, query *DeleteDocReq) *DocsAffected {
	var resp DocsAffected
	delete(index+"/", query, &resp)
	return &resp
}
