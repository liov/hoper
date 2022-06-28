package main

import (
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"gorm.io/gorm"
	"tools/timepill"
	"tools/timepill/model"
)

func main() {
	defer initialize.Start(&timepill.Conf, &timepill.Dao)()
	var user model.User
	err := timepill.Dao.Hoper.Where(`user_id = ?`, 100196722).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error(err)
	}
	timepill.RecordUserDiaries(&user)
}
