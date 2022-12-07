package dao

import (
	"context"
	initpostgres "github.com/liov/hoper/server/go/lib/initialize/db/postgres"
	initredis "github.com/liov/hoper/server/go/lib/initialize/redis"
	"gorm.io/gorm"
)

type dao3 struct {
	Hoper initpostgres.DB
	Redis initredis.Redis
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
