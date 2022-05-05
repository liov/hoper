package request

type SortType int

const (
	SortTypePlacement SortType = iota
	SortTypeASC
	SortTypeDESC
)

type ListReq struct {
	PageReq
	SortReq
}

type PageReq struct {
	PageNo   int `json:"pageNo"`
	PageSize int `json:"pageSize"`
}

type SortReq struct {
	SortField string   `json:"sortField"`
	SortType  SortType `json:"sortType"`
}

type DateReq = RangeReq

type RangeReq struct {
	RangeField string `json:"dateField"`
	RangeStart any    `json:"dateStart"`
	RangeEnd   any    `json:"dateEnd"`
	Include    bool   `json:"include"`
}
