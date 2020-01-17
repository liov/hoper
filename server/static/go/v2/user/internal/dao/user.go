package dao

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	modelconst "github.com/liov/hoper/go/v2/user/model"
	"github.com/liov/hoper/go/v2/utils/array"
	"github.com/liov/hoper/go/v2/utils/log"
)

type UserDao struct{}

func (*UserDao) ExitByEmailORPhone(mail, phone string) (bool, error) {
	var err error
	var count int
	if mail != "" {
		err = Dao.GORMDB.Table("user").Where(`mail = ?`, mail).Count(&count).Error
	} else {
		err = Dao.GORMDB.Table("user").Where(`phone = ?`, phone).Count(&count).Error
	}
	if err != nil {
		log.Error("UserDao.ExitByEmailORPhone: ", err)
		return true, err
	}
	return count == 1, nil
}

func (d *UserDao) GetByEmailORPhone(email, phone string, db *gorm.DB) (*model.User, error) {
	if db == nil {
		db = Dao.GORMDB
	}
	var user model.User
	var err error
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

func (*UserDao) Creat(user *model.User, db *gorm.DB) error {
	if db == nil {
		db = Dao.GORMDB
	}
	if err := db.Create(user).Error; err != nil {
		log.Error("UserDao.Creat: ", err)
		return err
	}
	return nil
}

func (*UserDao) GetByPrimaryKey(id uint64, db *gorm.DB) (*model.User, error) {
	if db == nil {
		db = Dao.GORMDB
	}
	var user model.User
	if err := db.First(&user, id).Error; err != nil {
		log.Error("UserDao.GetByPrimaryKey: ", err)
		return nil, err
	}
	return &user, nil
}

func (*UserDao) SaveResumes(userId uint64, resumes []*model.Resume, originalIds []uint64, device *model.UserDeviceInfo, db *gorm.DB) error {
	if db == nil {
		db = Dao.GORMDB
	}
	var err error
	var actionLog model.UserActionLog
	actionLog.CreatedAt = time.Now().Format(time.RFC3339Nano)
	actionLog.UserId = userId
	actionLog.UserDeviceInfo = device
	actionLog.Action = modelconst.EditResume
	tableName := db.NewScope(&model.Resume{}).TableName() + "."

	var editIds []uint64

	for i := range resumes {
		resumes[i].UserId = userId
		resumes[i].Status = 1
		if resumes[i].Id != 0 {
			err = db.Save(resumes[i]).Error
			actionLog.Action = modelconst.CreateResume
			actionLog.LastValue, _ = json.Marshal(resumes[i])
			editIds = append(editIds, resumes[i].Id)
		} else {
			err = db.Create(resumes[i]).Error
			actionLog.Action = modelconst.EditResume
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
		differenceIds = array.Difference(editIds, originalIds)
	}
	db.Model(&model.Resume{}).Where("id in (?)", differenceIds).Update("status", 0)
	for _, id := range differenceIds {
		actionLog.Id = 0
		actionLog.Action = modelconst.DELETEResume
		actionLog.RelatedId = tableName + strconv.FormatUint(id, 10)
		if err = db.Create(&actionLog).Error; err != nil {
			log.Error(err)
		}
	}
	return nil
}

func (*UserDao) ActionLog(log *model.UserActionLog, db *gorm.DB) error {
	err := db.Create(&log).Error
	if err != nil {
		return err
	}
	return nil
}

func (*UserDao) ResumesIds(userId uint64, db *gorm.DB) ([]uint64, error) {
	if db == nil {
		db = Dao.GORMDB
	}
	var resumeIds []uint64
	err := db.Model(new(model.Resume)).Where("user_id = ? AND status > 0", userId).Pluck("id", &resumeIds).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return resumeIds, nil
}
