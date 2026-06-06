package db

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	sqlx "github.com/hopeio/gox/database/sql"
	_ "github.com/hopeio/gox/database/sql/gorm"
	"github.com/hopeio/gox/log"
	"github.com/hopeio/gox/slices"
	puser "github.com/liov/hoper/server/go/protobuf/user"
	"github.com/liov/hoper/server/go/user/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserDao struct {
	*gorm.DB
}

func GetUserDao(db *gorm.DB) *UserDao {
	return &UserDao{db}
}

func (d *UserDao) GetByEmailOrPhone(ctx context.Context, mail string, countryCallingCode string, phone string, fields ...string) (*puser.User, error) {

	var u puser.User
	var err error
	db := d.DB
	if len(fields) > 0 {
		db = d.Table(model.TableNameUser).Select(fields)
	}
	if mail != "" {
		err = db.Where("mail = ? AND status != ?"+sqlx.WithNotDeleted, mail, puser.UserStatusDeleted).First(&u).Error
	} else if countryCallingCode != "" && phone != "" {
		err = db.Where("country_calling_code = ? AND phone = ? AND status != ?"+sqlx.WithNotDeleted, countryCallingCode, phone, puser.UserStatusDeleted).First(&u).Error
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (d *UserDao) Create(ctx context.Context, user *puser.User) error {
	return d.Table(model.TableNameUser).Create(user).Error
}

func (d *UserDao) GetByPrimaryKey(ctx context.Context, id uint64) (*puser.User, error) {

	var user puser.User
	if err := d.Table(model.TableNameUser).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *UserDao) SaveResumes(ctx context.Context, userId uint64, resumes []*puser.Resume, originalIds []uint64, device *puser.AccessDevice) error {

	if len(resumes) == 0 {
		return nil
	}
	var err error
	var actionLog puser.ActionLog
	//actionLog.CreatedAt = timestamp.New(time.Now())
	actionLog.UserId = userId
	actionLog.Device = device
	actionLog.Action = puser.ActionEditResume
	tableName := model.TableNameResume + "."

	var editIds []uint64

	for i := range resumes {
		resumes[i].UserId = userId
		resumes[i].Status = 1
		if resumes[i].Basic.Id != 0 {
			err = d.Save(resumes[i]).Error
			actionLog.Action = puser.ActionCreateResume
			actionLog.LastValue, _ = json.Marshal(resumes[i])
			editIds = append(editIds, resumes[i].Basic.Id)
		} else {
			err = d.DB.Create(resumes[i]).Error
			actionLog.Action = puser.ActionEditResume
		}
		if err != nil {
			return err
		}
		actionLog.Id = 0
		actionLog.RelatedId = tableName + strconv.FormatUint(resumes[i].Basic.Id, 10)
		if err = d.Table(model.TableNameActionLog).Create(&actionLog).Error; err != nil {
			log.Error(err)
		}

	}

	//差集
	var differenceIds []uint64
	if len(originalIds) > 0 {
		differenceIds = slices.DifferenceSet(editIds, originalIds)
	}
	d.Model(&puser.Resume{}).Where("id in (?)", differenceIds).Update("status", 0)
	for _, id := range differenceIds {
		actionLog.Id = 0
		actionLog.Action = puser.ActionDeleteResume
		actionLog.RelatedId = tableName + strconv.FormatUint(id, 10)
		if err = d.Table(model.TableNameActionLog).Create(&actionLog).Error; err != nil {
			return err
		}
	}
	return nil
}

func (d *UserDao) ActionLog(ctx context.Context, log *puser.ActionLog) error {

	err := d.Table(model.TableNameActionLog).Create(&log).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *UserDao) ResumesIds(ctx context.Context, userId uint64) ([]uint64, error) {

	var resumeIds []uint64
	err := d.Table(model.TableNameResume).Where("user_id = ? AND status > 0", userId).Pluck("id", &resumeIds).Error
	if err != nil {
		return nil, err
	}
	return resumeIds, nil
}

func (d *UserDao) GetBaseListDB(ctx context.Context, ids []uint64, pageNo, pageSize int) (int64, []*puser.UserBase, error) {

	db := d.Table(model.TableNameUser)
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
	var users []*puser.UserBase
	err := db.Clauses(clauses...).Scan(&users).Error
	if err != nil {
		return 0, nil, err
	}

	if len(ids) > 0 {
		count = int64(len(users))
	}
	return count, users, nil
}

func (d *UserDao) FollowExistsDB(ctx context.Context, id, followId uint64) (bool, error) {
	sql := `SELECT EXISTS(SELECT * FROM "` + model.TableNameFollow + `"
WHERE user_id = ?  AND follow_id = ?` + sqlx.WithNotDeleted + ` LIMIT 1)`
	var exists bool
	err := d.Raw(sql, id, followId).Scan(&exists).Error
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (d *UserDao) Active(ctx context.Context, u *puser.User) error {
	return d.Model(u).Updates(map[string]interface{}{"activated_at": time.Now(), "status": puser.UserStatusActivated}).Error
}

func (d *UserDao) Update(ctx context.Context, req *puser.EditReq) error {
	return d.Table(model.TableNameUser).Where(`id = ?`, req.Id).UpdateColumns(req.Detail).Error
}

func (d *UserDao) UserInfoByAccount(ctx context.Context, mail, countryCallingCode, phone string) (*puser.User, error) {
	var user puser.User
	var sql string
	if mail != "" {
		sql = "mail = ?"
	}else {
		sql = "country_calling_code = ? AND phone = ?"
	}
	sql += ` AND status != ?`+sqlx.WithNotDeleted
	db := d.Table(model.TableNameUser)
	if mail != "" {
		db = db.Where(sql, mail, puser.UserStatusDeleted)
	} else if countryCallingCode != "" && phone != "" {
		db = db.Where(sql, countryCallingCode, phone, puser.UserStatusDeleted)
	}
	return &user, db.First(&user).Error
}
