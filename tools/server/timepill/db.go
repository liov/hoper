package timepill

import (
	"context"
	"fmt"
	"github.com/actliboy/hoper/server/go/lib/utils/dao/db/gorm/clause"
	"github.com/actliboy/hoper/server/go/lib/utils/def/request"
	"gorm.io/gorm"
)

const (
	DiaryTableName = "diary"
)

type ListReq struct {
	request.ListReq
	request.RangeReq
}

func (dao *DBDao) ListDB(req *ListReq) ([]*Diary, error) {
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

type DBDao struct {
	ctx   context.Context
	Hoper *gorm.DB
}

func (dao *DBDao) MaxDiaryId() (int, error) {
	var maxId int

	err := dao.Hoper.Table(DiaryTableName).Select("MAX(id)").Scan(&maxId).Error
	if err != nil {
		return 0, err
	}
	return maxId, nil
}
