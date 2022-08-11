package dao

import (
	"context"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize/db"
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

func (d *dao) CreateVideo(video *Video) error {
	return d.db.Create(video).Error
}
