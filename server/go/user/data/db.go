package data

import (
	"encoding/json"
	dbi "github.com/hopeio/cherry/utils/dao/db"
	"github.com/liov/hoper/server/go/user/confdao"
	"gorm.io/gorm/clause"
	"strconv"
	"time"

	"github.com/hopeio/cherry/protobuf/errorcode"
	"github.com/hopeio/cherry/utils/log"
	"github.com/hopeio/cherry/utils/slices"
	"github.com/liov/hoper/server/go/protobuf/user"
	"github.com/liov/hoper/server/go/user/model"
	"gorm.io/gorm"
)

func DBNotNil(db **gorm.DB) {
	if *db == nil {
		*db = confdao.Dao.GORMDB.DB
	}
}

func (d *userDao) ExitsCheck(db *gorm.DB, field, value string) (bool, error) {
	DBNotNil(&db)
	ctxi := d
	sql := `SELECT EXISTS(SELECT id FROM "` + model.UserTableName + `" WHERE ` + dbi.NotDeleted + ` AND `
	var exists bool
	err := db.Raw(sql+field+` = ? AND status != ?  LIMIT 1)`, value, user.UserStatusDeleted).Scan(&exists).Error

	if err != nil {
		return true, ctxi.ErrorLog(errorcode.DBError, err, "ExitsCheck")
	}
	return exists, nil
}

func (d *userDao) GetByEmailORPhone(db *gorm.DB, email, phone string, fields ...string) (*user.User, error) {
	DBNotNil(&db)
	ctxi := d
	var u user.User
	var err error
	if len(fields) > 0 {
		db = db.Table(model.UserTableName).Select(fields)
	}
	if email != "" {
		err = db.Where("email = ? AND status != ?"+dbi.WithNotDeleted, email, user.UserStatusDeleted).Find(&u).Error
	} else {
		err = db.Where("phone = ? AND status != ?"+dbi.WithNotDeleted, phone, user.UserStatusDeleted).Find(&u).Error
	}
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError, err, "GetByEmailORPhone")
	}
	return &u, nil
}

func (*userDao) Creat(db *gorm.DB, user *user.User) error {
	DBNotNil(&db)
	if err := db.Table(model.UserTableName).Create(user).Error; err != nil {
		log.Error("UserDao.Creat: ", err)
		return err
	}
	return nil
}

func (d *userDao) GetByPrimaryKey(db *gorm.DB, id uint64) (*user.User, error) {
	DBNotNil(&db)
	ctxi := d
	var user user.User
	if err := db.Table(model.UserTableName).First(&user, id).Error; err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError, err, "GetByPrimaryKey")
	}
	return &user, nil
}

func (d *userDao) SaveResumes(db *gorm.DB, userId uint64, resumes []*user.Resume, originalIds []uint64, device *user.UserDeviceInfo) error {
	DBNotNil(&db)
	ctxi := d
	if len(resumes) == 0 {
		return nil
	}
	var err error
	var actionLog user.UserActionLog
	actionLog.CreatedAt = time.Now().Format(time.RFC3339Nano)
	actionLog.UserId = userId
	actionLog.DeviceInfo = device
	actionLog.Action = user.ActionEditResume
	tableName := model.ResumeTableName + "."

	var editIds []uint64

	for i := range resumes {
		resumes[i].UserId = userId
		resumes[i].Status = 1
		if resumes[i].Id != 0 {
			err = db.Save(resumes[i]).Error
			actionLog.Action = user.ActionCreateResume
			actionLog.LastValue, _ = json.Marshal(resumes[i])
			editIds = append(editIds, resumes[i].Id)
		} else {
			err = db.Create(resumes[i]).Error
			actionLog.Action = user.ActionEditResume
		}
		if err != nil {
			return ctxi.ErrorLog(errorcode.DBError, err, "Save")
		}
		actionLog.Id = 0
		actionLog.RelatedId = tableName + strconv.FormatUint(resumes[i].Id, 10)
		if err = db.Table(model.UserActionLogTableName).Create(&actionLog).Error; err != nil {
			log.Error(err)
		}

	}

	//差集
	var differenceIds []uint64
	if len(originalIds) > 0 {
		differenceIds = slices.DifferenceSet(editIds, originalIds)
	}
	db.Model(&user.Resume{}).Where("id in (?)", differenceIds).Update("status", 0)
	for _, id := range differenceIds {
		actionLog.Id = 0
		actionLog.Action = user.ActionDELETEResume
		actionLog.RelatedId = tableName + strconv.FormatUint(id, 10)
		if err = db.Table(model.UserActionLogTableName).Create(&actionLog).Error; err != nil {
			log.Error(err)
		}
	}
	return nil
}

func (d *userDao) ActionLog(db *gorm.DB, log *user.UserActionLog) error {
	ctxi := d
	err := db.Table(model.UserActionLogTableName).Create(&log).Error
	if err != nil {
		return ctxi.ErrorLog(errorcode.DBError, err, "ActionLog")
	}
	return nil
}

func (d *userDao) ResumesIds(db *gorm.DB, userId uint64) ([]uint64, error) {
	DBNotNil(&db)
	ctxi := d
	var resumeIds []uint64
	err := db.Table(model.ResumeTableName).Where("user_id = ? AND status > 0", userId).Pluck("id", &resumeIds).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, ctxi.ErrorLog(errorcode.DBError, err, "ResumesIds")
	}
	return resumeIds, nil
}

func (d *userDao) GetBaseListDB(db *gorm.DB, ids []uint64, pageNo, pageSize int) (int64, []*user.UserBaseInfo, error) {
	ctxi := d

	db = db.Table(model.UserTableName)
	var count int64
	if len(ids) > 0 {
		db = db.Where("id IN (?)", ids)
	} else {
		err := db.Count(&count).Error
		if err != nil {
			return 0, nil, ctxi.ErrorLog(errorcode.DBError, err, "Count")
		}
	}

	var clauses []clause.Expression
	if pageNo != 0 && pageSize != 0 {
		clauses = append(clauses, clause.Limit{Offset: (pageNo - 1) * pageSize, Limit: &pageSize})
	}
	var users []*user.UserBaseInfo
	err := db.Clauses(clauses...).Scan(&users).Error
	if err != nil {
		return 0, nil, ctxi.ErrorLog(errorcode.DBError, err, "Scan")
	}

	if len(ids) > 0 {
		count = int64(len(users))
	}
	return count, users, nil
}

func (d *userDao) FollowExistsDB(db *gorm.DB, id, followId uint64) (bool, error) {
	ctxi := d
	sql := `SELECT EXISTS(SELECT * FROM "` + model.FollowTableName + `" 
WHERE user_id = ?  AND follow_id = ?` + dbi.WithNotDeleted + ` LIMIT 1)`
	var exists bool
	err := db.Raw(sql, id, followId).Scan(&exists).Error
	if err != nil {
		return false, ctxi.ErrorLog(errorcode.DBError, err, "ExistsByAuth")
	}
	return exists, nil
}
