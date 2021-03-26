package dao

import (
	"database/sql"
	"github.com/liov/hoper/go/v2/protobuf/utils/errorcode"
	"github.com/liov/hoper/go/v2/upload/model"
	"gorm.io/gorm"
)

func (d *uploadDao) UploadDB(db *gorm.DB,md5,size string) (*model.FileUploadInfo,error){
	var upload *model.FileUploadInfo
	raw := `SELECT id FROM upload WHERE md5 = ? AND size = ? LIMIT 1`
	var id uint64
	err:=db.Raw(raw,md5,size).Row().Scan(&id)
	if err != nil && err !=sql.ErrNoRows{
		return nil,d.ctxi.ErrorLog(errorcode.DBError,err,"FileExistsDB")
	}
	return upload,nil
}
