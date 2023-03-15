package _type

import (
	"github.com/liov/hoper/server/go/lib/utils/def/request"
	"golang.org/x/exp/constraints"
	"time"
)

type ListReq[T Ordered] struct {
	request.PageSortReq
	*request.RangeReq[T]
}

func NewListReq[T Ordered](pageNo, pageSize int) *ListReq[T] {
	return &ListReq[T]{
		PageSortReq: request.PageSortReq{
			PageReq: request.PageReq{
				PageNo:   pageNo,
				PageSize: pageSize,
			},
		},
	}
}

func (req *ListReq[T]) WithSort(field string, typ request.SortType) *ListReq[T] {
	req.SortReq = &request.SortReq{
		SortField: field,
		SortType:  typ,
	}
	return req
}

func (req *ListReq[T]) WithRange(field string, start, end T, include bool) *ListReq[T] {
	req.RangeReq = &request.RangeReq[T]{
		RangeField: field,
		RangeStart: start,
		RangeEnd:   end,
		Include:    include,
	}
	return req
}

type Ordered interface {
	constraints.Ordered | time.Time
}
