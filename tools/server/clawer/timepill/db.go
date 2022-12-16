package timepill

import (
	"context"
	"fmt"
	clausei "github.com/liov/hoper/server/go/lib_v2/utils/dao/db/gorm/clause"
	_type "github.com/liov/hoper/server/go/lib_v2/utils/dao/db/type"
	"gorm.io/gorm"

	"tools/clawer/timepill/model"
)

func (dao *DBDao) List(req *_type.ListReq[int]) ([]*model.Diary, error) {
	return clausei.List[*model.Diary, int](Dao.Hoper.DB, req)
}

func CreateTable() {
	fmt.Println(Dao.Hoper.Migrator().CreateTable(&model.Badge{}, &model.User{}, &model.Diary{}, &model.NoteBook{}, &model.Comment{}, &model.Face{}))
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

	err := dao.Hoper.Table(model.DiaryTableName).Select("COALESCE(MAX(id),10000)").Scan(&maxId).Error
	if err != nil {
		return 0, err
	}
	return maxId, nil
}
