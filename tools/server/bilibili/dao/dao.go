package dao

import (
	"context"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize/db"
	"github.com/actliboy/hoper/server/go/lib/utils/dao/db/gorm/postgres"
	"gorm.io/gorm"
)

type dao3 struct {
	Hoper db.DB
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

func (d *dao) CreateView(view *View) error {
	return d.db.Create(view).Error
}

func (d *dao) ViewExists(aid int) (bool, error) {
	return postgres.Exists(d.db, TableNameView, "aid", aid, false)
}

func (d *dao) CreateVideo(video *Video) error {
	return d.db.Create(video).Error
}

func (d *dao) VideoExists(aid, cid int) (bool, error) {
	return postgres.ExistsBySQL(d.db, `SELECT EXISTS(SELECT * FROM `+TableNameVideo+` WHERE aid = ? AND cid = ? LIMIT 1)`, aid, cid)
}
