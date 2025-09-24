package db

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/hopeio/context/httpctx"
	sqlx "github.com/hopeio/gox/database/sql"
	gormx "github.com/hopeio/gox/database/sql/gorm"
	_ "github.com/hopeio/gox/database/sql/gorm/serializer"
	"github.com/hopeio/gox/log"
	"github.com/hopeio/gox/slices"
	"github.com/hopeio/gox/validation/validator"
	puser "github.com/liov/hoper/server/go/protobuf/user"
	"github.com/liov/hoper/server/go/user/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserDao struct {
	*httpctx.Context
	*gorm.DB
}

func GetUserDao(ctx *httpctx.Context, db *gorm.DB) *UserDao {
	if ctx == nil {
		log.Fatal("ctx can't nil")
	}
	return &UserDao{ctx, gormx.NewTraceDB(db, ctx.Base(), ctx.TraceID())}
}

func (d *UserDao) GetByNameOrEmailOrPhone(name, email, phone string) (*model.User, error) {

	var u model.User
	var err error
	err = d.Where("(name = ? OR mail = ? OR phone = ?) AND status != ?", name, email, phone, puser.UserStatusDeleted).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (d *UserDao) GetByEmailOrPhone(input string, fields ...string) (*puser.User, error) {

	var u puser.User
	var err error
	db := d.DB
	if len(fields) > 0 {
		db = d.Table(model.TableNameUser).Select(fields)
	}
	if strings.Contains(input, "@") {
		err = db.Where("mail = ? AND status != ?"+sqlx.WithNotDeleted, input, puser.UserStatusDeleted).First(&u).Error
	} else {
		err = db.Where("phone = ? AND status != ?"+sqlx.WithNotDeleted, input, puser.UserStatusDeleted).First(&u).Error
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (d *UserDao) Creat(user *puser.User) error {
	if err := d.Table(model.TableNameUser).Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (d *UserDao) GetByPrimaryKey(id uint64) (*puser.User, error) {

	var user puser.User
	if err := d.Table(model.TableNameUser).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *UserDao) SaveResumes(userId uint64, resumes []*puser.Resume, originalIds []uint64, device *puser.AccessDevice) error {

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
		if resumes[i].Id != 0 {
			err = d.Save(resumes[i]).Error
			actionLog.Action = puser.ActionCreateResume
			actionLog.LastValue, _ = json.Marshal(resumes[i])
			editIds = append(editIds, resumes[i].Id)
		} else {
			err = d.Create(resumes[i]).Error
			actionLog.Action = puser.ActionEditResume
		}
		if err != nil {
			return err
		}
		actionLog.Id = 0
		actionLog.RelatedId = tableName + strconv.FormatUint(resumes[i].Id, 10)
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

func (d *UserDao) ActionLog(log *puser.ActionLog) error {

	err := d.Table(model.TableNameActionLog).Create(&log).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *UserDao) ResumesIds(userId uint64) ([]uint64, error) {

	var resumeIds []uint64
	err := d.Table(model.TableNameResume).Where("user_id = ? AND status > 0", userId).Pluck("id", &resumeIds).Error
	if err != nil {
		return nil, err
	}
	return resumeIds, nil
}

func (d *UserDao) GetBaseListDB(ids []uint64, pageNo, pageSize int) (int64, []*puser.UserBase, error) {

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

func (d *UserDao) FollowExistsDB(id, followId uint64) (bool, error) {
	sql := `SELECT EXISTS(SELECT * FROM "` + model.TableNameFollow + `" 
WHERE user_id = ?  AND follow_id = ?` + sqlx.WithNotDeleted + ` LIMIT 1)`
	var exists bool
	err := d.Raw(sql, id, followId).Scan(&exists).Error
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (d *UserDao) Active(u *puser.User) error {
	return d.Model(u).Updates(map[string]interface{}{"activated_at": time.Now(), "status": puser.UserStatusActivated}).Error
}

func (d *UserDao) Update(req *puser.EditReq) error {
	return d.Table(model.TableNameUser).Where(`id = ?`, req.Id).UpdateColumns(req.Detail).Error
}

func (d *UserDao) UserInfoByAccount(account string) (*puser.User, error) {
	var user puser.User
	var sql string
	switch validator.PhoneOrMail(account) {
	case validator.Mail:
		sql = "mail = ?"
	case validator.Phone:
		sql = "phone = ?"
	default:
		sql = "account = ?"
	}
	return &user, d.Table(model.TableNameUser).
		Where(sql+` AND status != ?`+sqlx.WithNotDeleted, account, puser.UserStatusDeleted).First(&user).Error
}
