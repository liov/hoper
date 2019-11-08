package dao

import (
	"time"

	"github.com/liov/hoper/go/v2/protobuf/user/model"
	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/time2"
)

type UserDao struct {}

func (*UserDao) ExitByEmailORPhone(email,phone string) (bool,error)  {
	var err error
	var count int
	if email != "" {
		err = Dao.DB.DB().QueryRow(`SELECT EXISTS(select  1 FROM user WHERE email =  ?)`,email).Scan(&count)
	}else {
		err = Dao.DB.DB().QueryRow(`SELECT EXISTS(select  1 FROM user WHERE phone =  ?)`,phone).Scan(&count)
	}
	if  err != nil {
		log.Error("UserDao.ExitByEmailORPhone: ",err)
		return false,err
	}
	return count == 1,nil
}


func (*UserDao) GetByEmailORPhone(email,phone string) (*model.User,error) {
	var user model.User
	var err error
	if email != "" {
		err = Dao.DB.Where("email = ?", email).Find(&user).Error
	}else {
		err = Dao.DB.Where("phone = ?", phone).Find(&user).Error
	}
	if  err != nil {
		log.Error("UserDao.GetByEmailORPhone: ",err)
		return nil,err
	}
	return &user,nil
}

func (*UserDao) Creat(user *model.User) error {
	defer time2.TimeCost(time.Now())
	res,err := Dao.DB.DB().Exec(`INSERT INTO user 
    (name, password, email, phone, gender, role, avatar_url, created_at, updated_at, status)
     VALUES (?,?,?,?,?,?,?,?,?,?)`,user.Name,user.Password,user.Email,user.Phone,user.Role,
     user.AvatarURL,user.CreatedAt,user.UpdatedAt,user.Status)
	if  err != nil {
		log.Error("UserDao.Creat: ",err)
		return err
	}
	id,err := res.LastInsertId()
	user.Id = uint64(id)
	return nil
}