package timepill

import (
	"github.com/actliboy/hoper/server/go/lib/utils/dao/db/gorm/clause"
	"github.com/actliboy/hoper/server/go/lib/utils/def/request"
)

type ListReq struct {
	request.ListReq
	request.RangeReq
}

func List(req *ListReq) ([]*Diary, error) {
	var diaries []*Diary

	clauses := append((*clause.ListReq)(&req.ListReq).Clause(), (*clause.RangeReq)(&req.RangeReq).Clause())
	err := Dao.Hoper.Clauses(clauses...).Find(&diaries).Error
	if err != nil {
		return nil, err
	}
	return diaries, nil
}
