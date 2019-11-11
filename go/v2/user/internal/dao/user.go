package dao

import (
	"time"

	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/time2"
)

type UserDao struct{}

func (*UserDao) ExitByEmailORPhone(email, phone string) (bool, error) {
	var err error
	var count int
	if email != "" {
		err = Dao.DB.QueryRow(`SELECT EXISTS(select  1 FROM user WHERE email =  ?)`, email).Scan(&count)
	} else {
		err = Dao.DB.QueryRow(`SELECT EXISTS(select  1 FROM user WHERE phone =  ?)`, phone).Scan(&count)
	}
	if err != nil {
		log.Error("UserDao.ExitByEmailORPhone: ", err)
		return false, err
	}
	return count == 1, nil
}

func (*UserDao) GetByEmailORPhone(email, phone string) (*model.User, error) {
	var user model.User
	var err error
	if email != "" {
		err = Dao.GORMDB.Where("email = ?", email).Find(&user).Error
	} else {
		err = Dao.GORMDB.Where("phone = ?", phone).Find(&user).Error
	}
	if err != nil {
		log.Error("UserDao.GetByEmailORPhone: ", err)
		return nil, err
	}
	return &user, nil
}

func (*UserDao) Creat(user *model.User) error {
	defer time2.TimeCost(time.Now())
	res,err :=Dao.DB.Exec(`INSERT INTO user 
    ( name, password, email, phone, gender, avatar_url, role, created_at, updated_at, status, last_active_at, score, follow_count, followed_count, article_count, moment_count, diary_book_count, diary_count, comment_count) VALUES 
    (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		user.Name,user.Password,user.Email,user.Phone,user.Gender,user.AvatarURL,user.Role,user.CreatedAt,
		user.UpdatedAt,0,user.CreatedAt,0,0,0,0,0,0,0,0)
	if err!=nil {
		log.Error("UserDao.Creat: ", err)
		return err
	}
	id,err := res.LastInsertId()
	if err!=nil {
		log.Error("UserDao.LastInsertId: ", err)
		return err
	}
	user.Id = uint64(id)
	return nil
}
