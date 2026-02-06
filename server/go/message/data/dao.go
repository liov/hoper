package data

import (
	"github.com/hopeio/gox/context/httpctx"
	"log"
)

type chatDao struct {
	context.Context
}

func GetDao(ctx context.Context) *chatDao {
	if ctx == nil {
		log.Fatal("ctx can't nil")
	}
	return &chatDao{ctx}
}
