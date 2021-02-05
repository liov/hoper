package dao

import (
	"encoding/json"
	"strconv"
	"time"

	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/slices"
	"gorm.io/gorm"
)

type UserDao struct{}

func DBNotNil(db **gorm.DB) {
	if *db == nil {
		*db = Dao.GORMDB
	}
}

func (*UserDao) ExitByEmailORPhone(ctx *model.Ctx,db *gorm.DB, mail, phone string) (bool, error) {
	DBNotNil(&db)
	var err error
	var count int64
	if mail != "" {
		err = db.Table("user").Where(`mail = ?`, mail).Count(&count).Error
	} else {
		err = db.Table("user").Where(`phone = ?`, phone).Count(&count).Error
	}
	if err != nil {
		log.Error("UserDao.ExitByEmailORPhone: ", err)
		return true, err
	}
	return count == 1, nil
}

func (d *UserDao) GetByEmailORPhone(ctx *model.Ctx,db *gorm.DB, email, phone string, fields ...string) (*model.User, error) {
	DBNotNil(&db)
	var user model.User
	var err error
	if len(fields) > 0 {
		db = db.Select(fields)
	}
	if email != "" {
		err = db.Where("email = ?", email).Find(&user).Error
	} else {
		err = db.Where("phone = ?", phone).Find(&user).Error
	}
	if err != nil {
		log.Error("UserDao.GetByEmailORPhone: ", err)
		return nil, err
	}
	return &user, nil
}

func (*UserDao) Creat(ctx *model.Ctx,db *gorm.DB, user *model.User) error {
	DBNotNil(&db)
	if err := db.Create(user).Error; err != nil {
		log.Error("UserDao.Creat: ", err)
		return err
	}
	return nil
}

func (*UserDao) GetByPrimaryKey(ctx *model.Ctx,db *gorm.DB, id uint64) (*model.User, error) {
	DBNotNil(&db)
	var user model.User
	if err := db.First(&user, id).Error; err != nil {
		log.Error("UserDao.GetByPrimaryKey: ", err)
		return nil, err
	}
	return &user, nil
}

func (*UserDao) SaveResumes(ctx *model.Ctx,db *gorm.DB, userId uint64, resumes []*model.Resume, originalIds []uint64, device *model.UserDeviceInfo) error {
	DBNotNil(&db)
	if len(resumes) == 0 {
		return nil
	}
	var err error
	var actionLog model.UserActionLog
	actionLog.CreatedAt = time.Now().Format(time.RFC3339Nano)
	actionLog.UserId = userId
	actionLog.DeviceInfo = device
	actionLog.Action = model.ActionEditResume
	tableName := resumes[0].TableName() + "."

	var editIds []uint64

	for i := range resumes {
		resumes[i].UserId = userId
		resumes[i].Status = 1
		if resumes[i].Id != 0 {
			err = db.Save(resumes[i]).Error
			actionLog.Action = model.ActionCreateResume
			actionLog.LastValue, _ = json.Marshal(resumes[i])
			editIds = append(editIds, resumes[i].Id)
		} else {
			err = db.Create(resumes[i]).Error
			actionLog.Action = model.ActionEditResume
		}
		if err != nil {
			return err
		}
		actionLog.Id = 0
		actionLog.RelatedId = tableName + strconv.FormatUint(resumes[i].Id, 10)
		if err = db.Create(&actionLog).Error; err != nil {
			log.Error(err)
		}

	}

	//差集
	var differenceIds []uint64
	if len(originalIds) > 0 {
		differenceIds = slices.Difference(editIds, originalIds)
	}
	db.Model(&model.Resume{}).Where("id in (?)", differenceIds).Update("status", 0)
	for _, id := range differenceIds {
		actionLog.Id = 0
		actionLog.Action = model.ActionDELETEResume
		actionLog.RelatedId = tableName + strconv.FormatUint(id, 10)
		if err = db.Create(&actionLog).Error; err != nil {
			log.Error(err)
		}
	}
	return nil
}

func (*UserDao) ActionLog(ctx *model.Ctx,db *gorm.DB, log *model.UserActionLog) error {
	err := db.Create(&log).Error
	if err != nil {
		return err
	}
	return nil
}

func (*UserDao) ResumesIds(ctx *model.Ctx,db *gorm.DB, userId uint64) ([]uint64, error) {
	DBNotNil(&db)
	var resumeIds []uint64
	err := db.Model(new(model.Resume)).Where("user_id = ? AND status > 0", userId).Pluck("id", &resumeIds).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return resumeIds, nil
}
