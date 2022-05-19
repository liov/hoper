package timepill

import (
	"fmt"
	"github.com/actliboy/hoper/server/go/lib/utils/dao/db/gorm/clause"
	"github.com/actliboy/hoper/server/go/lib/utils/def/request"
)

type ListReq struct {
	request.ListReq
	request.RangeReq
}

func List(req *ListReq) ([]*Diary, error) {
	var diaries []*Diary

	clauses := append((*clausei.ListReq)(&req.ListReq).Clause(), (*clausei.RangeReq)(&req.RangeReq).Clause())
	err := Dao.Hoper.Clauses(clauses...).Find(&diaries).Error
	if err != nil {
		return nil, err
	}
	return diaries, nil
}

func CreateTable() {
	fmt.Println(Dao.Hoper.Migrator().CreateTable(&Badge{}, &User{}, &Diary{}, &NoteBook{}, &Comment{}))
}

func CreateBadgeTable() {
	fmt.Println(Dao.Hoper.Migrator().CreateTable(&Badge{}))
}

func CreateCommentTable() {
	fmt.Println(Dao.Hoper.Migrator().CreateTable(&Comment{}))
}
