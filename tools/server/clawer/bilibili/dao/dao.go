package dao

import (
	"context"
	initpostgres "github.com/liov/hoper/server/go/lib/initialize/db/postgres"
	"github.com/liov/hoper/server/go/lib/utils/dao/db/gorm/postgres"
	"time"

	"gorm.io/gorm"
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

func (d *dao) CreateView(view *View) error {
	return d.db.Create(view).Error
}

func (d *dao) ViewExists(aid int) (bool, error) {
	return postgres.Exists(d.db, TableNameView, "aid", aid, false)
}

func (d *dao) ViewInfo(aid int) (*View, error) {
	var view View
	err := d.db.Table(TableNameView).Where(`aid = ?`, aid).First(&view).Error
	if err != nil {
		return nil, err
	}
	return &view, nil
}

func (d *dao) ViewCreatedTime(aid int) (time.Time, error) {
	var t time.Time
	err := d.db.Table(TableNameView).Select(`created_at`).Where(`aid = ?`, aid).Scan(&t).Error
	if err != nil {
		return t, err
	}
	return t, nil
}

func (d *dao) LastCreated(tableName string) (time.Time, error) {
	var t time.Time
	err := d.db.Table(tableName).Select(`created_at`).Order(`created_at DESC`).Limit(1).Scan(&t).Error
	if err != nil {
		return t, err
	}
	return t, nil
}

func (d *dao) LastRecordAid() (int, error) {
	var id int
	err := d.db.Table(TableNameView).Select(`aid`).Order(`created_at DESC`).Limit(1).Scan(&id).Error
	if err != nil {
		return id, err
	}
	return id, nil
}

func (d *dao) CreateVideo(video *Video) error {
	return d.db.Create(video).Error
}

func (d *dao) VideoExists(cid int) (bool, error) {
	return postgres.ExistsBySQL(d.db, `SELECT EXISTS(SELECT * FROM `+TableNameVideo+` WHERE cid = ? LIMIT 1)`, cid)
}
