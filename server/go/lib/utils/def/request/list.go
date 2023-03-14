package request

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

type RangeReq struct {
	RangeField string      `json:"dateField,omitempty"`
	RangeStart interface{} `json:"dateStart,omitempty"`
	RangeEnd   interface{} `json:"dateEnd,omitempty"`
	Include    bool        `json:"include,omitempty"`
}
