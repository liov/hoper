package timepill

import (
	"context"
	"fmt"
	clausei "github.com/actliboy/hoper/server/go/lib/utils/generics/dao/db/gorm/clause"
	_type "github.com/actliboy/hoper/server/go/lib/utils/generics/dao/db/type"
	"gorm.io/gorm"

	"tools/timepill/model"
)

const (
	DiaryTableName = "diary"
	FaceTableName  = "face"
)

func (dao *DBDao) List(req *_type.ListReq[int]) ([]*model.Diary, error) {
	return clausei.List[*model.Diary, int](Dao.Hoper.DB, req)
}

func CreateTable() {
	fmt.Println(Dao.Hoper.Migrator().CreateTable(&model.Badge{}, &model.User{}, &model.Diary{}, &model.NoteBook{}, &model.Comment{}))
}

func CreateBadgeTable() {
	fmt.Println(Dao.Hoper.Migrator().CreateTable(&model.Badge{}))
}

func CreateCommentTable() {
	fmt.Println(Dao.Hoper.Migrator().CreateTable(&model.Comment{}))
}

func CreateFaceTable() {
	fmt.Println(Dao.Hoper.Migrator().CreateTable(&model.Face{}))
}

type DBDao struct {
	ctx   context.Context
	Hoper *gorm.DB
}

func (dao *DBDao) MaxDiaryId() (int, error) {
	var maxId int

	err := dao.Hoper.Table(FaceTableName).Select("MAX(id)").Scan(&maxId).Error
	if err != nil {
		return 0, err
	}
	return maxId, nil
}
