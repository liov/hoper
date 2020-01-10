package dao

import (
	"github.com/jinzhu/gorm"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	modelconst "github.com/liov/hoper/go/v2/user/internal/model"
	"github.com/liov/hoper/go/v2/utils/log"
)

type UserDao struct{}

func (*UserDao) ExitByEmailORPhone(email, phone string) (bool, error) {
	var err error
	var count int
	if email != "" {
		err = Dao.StdDB.QueryRow(`SELECT EXISTS(select  1 FROM user WHERE mail =  ?)`, email).Scan(&count)
	} else {
		err = Dao.StdDB.QueryRow(`SELECT EXISTS(select  1 FROM user WHERE phone =  ?)`, phone).Scan(&count)
	}
	if err != nil {
		log.Error("UserDao.ExitByEmailORPhone: ", err)
		return false, err
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

func (*UserDao) SaveResume(userId uint64, resumes []*model.Resume, device *model.UserDeviceInfo, db *gorm.DB) error {
	if db == nil {
		db = Dao.GORMDB
	}
	var err error
	var log model.UserActionLog
	log.UserId = userId
	log.UserDeviceInfo = device
	log.Action = modelconst.EditResume
	for _, resume := range resumes {
		resume.UserId = userId
		resume.Status = 1
		if resume.Id != 0 {
			err = db.Model(&resume).Updates(&resume).Error
		} else {
			err = db.Create(&resume).Error
		}
		if err != nil {
			return err
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
