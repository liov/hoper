package data

import (
	"encoding/json"
	dbi "github.com/hopeio/cherry/utils/dao/db"
	"github.com/hopeio/cherry/utils/log"
	"github.com/hopeio/cherry/utils/slices"
	"github.com/liov/hoper/server/go/protobuf/user"
	"github.com/liov/hoper/server/go/user/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
	"strings"
)

func (d *userDao) GetByNameOrEmailOrPhone(db *gorm.DB, name, email, phone string) (*model.User, error) {

	var u model.User
	var err error
	err = db.Where("(name = ? OR mail = ? OR phone = ?) AND status != ?", name, email, phone, user.UserStatusDeleted).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (d *userDao) GetByEmailOrPhone(db *gorm.DB, input string, fields ...string) (*user.User, error) {

	var u user.User
	var err error
	if len(fields) > 0 {
		db = db.Table(model.UserTableName).Select(fields)
	}
	if strings.Contains(input, "@") {
		err = db.Where("mail = ? AND status != ?"+dbi.WithNotDeleted, input, user.UserStatusDeleted).First(&u).Error
	} else {
		err = db.Where("phone = ? AND status != ?"+dbi.WithNotDeleted, input, user.UserStatusDeleted).First(&u).Error
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (*userDao) Creat(db *gorm.DB, user *user.User) error {
	if err := db.Table(model.UserTableName).Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (d *userDao) GetByPrimaryKey(db *gorm.DB, id uint64) (*user.User, error) {

	var user user.User
	if err := db.Table(model.UserTableName).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *userDao) SaveResumes(db *gorm.DB, userId uint64, resumes []*user.Resume, originalIds []uint64, device *user.UserDeviceInfo) error {

	if len(resumes) == 0 {
		return nil
	}
	var err error
	var actionLog user.UserActionLog
	//actionLog.CreatedAt = timestamp.New(time.Now())
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
			return err
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
			return err
		}
	}
	return nil
}

func (d *userDao) ActionLog(db *gorm.DB, log *user.UserActionLog) error {

	err := db.Table(model.UserActionLogTableName).Create(&log).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *userDao) ResumesIds(db *gorm.DB, userId uint64) ([]uint64, error) {

	var resumeIds []uint64
	err := db.Table(model.ResumeTableName).Where("user_id = ? AND status > 0", userId).Pluck("id", &resumeIds).Error
	if err != nil {
		return nil, err
	}
	return resumeIds, nil
}

func (d *userDao) GetBaseListDB(db *gorm.DB, ids []uint64, pageNo, pageSize int) (int64, []*user.UserBaseInfo, error) {

	db = db.Table(model.UserTableName)
	var count int64
	if len(ids) > 0 {
		db = db.Where("id IN (?)", ids)
	} else {
		err := db.Count(&count).Error
		if err != nil {
			return 0, nil, err
		}
	}

	var clauses []clause.Expression
	if pageNo != 0 && pageSize != 0 {
		clauses = append(clauses, clause.Limit{Offset: (pageNo - 1) * pageSize, Limit: &pageSize})
	}
	var users []*user.UserBaseInfo
	err := db.Clauses(clauses...).Scan(&users).Error
	if err != nil {
		return 0, nil, err
	}

	if len(ids) > 0 {
		count = int64(len(users))
	}
	return count, users, nil
}

func (d *userDao) FollowExistsDB(db *gorm.DB, id, followId uint64) (bool, error) {
	sql := `SELECT EXISTS(SELECT * FROM "` + model.FollowTableName + `" 
WHERE user_id = ?  AND follow_id = ?` + dbi.WithNotDeleted + ` LIMIT 1)`
	var exists bool
	err := db.Raw(sql, id, followId).Scan(&exists).Error
	if err != nil {
		return false, err
	}
	return exists, nil
}
