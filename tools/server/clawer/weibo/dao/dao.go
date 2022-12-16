package dao

import (
	"context"
	initpostgres "github.com/liov/hoper/server/go/lib/initialize/gormdb/postgres"
	"github.com/liov/hoper/server/go/lib/utils/dao/db/gorm/postgres"
	"gorm.io/gorm"
	"strconv"
)

type dao3 struct {
	Hoper initpostgres.DB
}

func (d dao3) Init() {
}

var Dao dao3

type dao struct {
	ctx context.Context
	db  *gorm.DB
}

func NewDao(ctx context.Context, db *gorm.DB) *dao {
	return &dao{ctx, db}
}

func WeiboExists(id string) (bool, error) {
	return postgres.ExistsBySQL(Dao.Hoper.DB, `SELECT EXISTS(SELECT * FROM `+TableNameWeibo+` WHERE id = `+id+` LIMIT 1)`)
}

func UserExists(id int) (bool, error) {
	return postgres.ExistsBySQL(Dao.Hoper.DB, `SELECT EXISTS(SELECT * FROM `+TableNameUser+` WHERE id = `+strconv.Itoa(id)+` LIMIT 1)`)
}
