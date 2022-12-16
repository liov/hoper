package request

import (
	"github.com/liov/hoper/server/go/lib_v2/utils/def/constraints"
	"time"
)

type SortType int

const (
	SortTypePlacement SortType = iota
	SortTypeASC
	SortTypeDESC
)

type PageSortReqInter interface {
	PageReqInter
	SortReqInter
}

type PageSortReq struct {
	PageReq
	*SortReq
}

type PageReqInter interface {
	PageNo() int
	PageSize() int
}

type PageReq struct {
	PageNo   int `json:"pageNo"`
	PageSize int `json:"pageSize"`
}

type SortReq struct {
	SortField string   `json:"sortField,omitempty"`
	SortType  SortType `json:"sortType,omitempty"`
}

func (receiver *SortReq) Column() string {
	return receiver.SortField
}

func (receiver *SortReq) Type() SortType {
	return receiver.SortType
}

type SortReqInter interface {
	Column() string
	Type() SortType
}

type DateRangeReq[T ~string | time.Time] RangeReq[T]

type RangeReq[T constraints.Ordered] struct {
	RangeField string `json:"dateField,omitempty"`
	RangeStart T      `json:"dateStart,omitempty"`
	RangeEnd   T      `json:"dateEnd,omitempty"`
	Include    bool   `json:"include,omitempty"`
}
