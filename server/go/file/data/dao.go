package data

import (
	"github.com/hopeio/gox/context/httpctx"
	"log"
)

type uploadDao struct {
	context.Context
}

func GetDao(ctx context.Context) *uploadDao {
	if ctx == nil {
		log.Fatal("ctx can't nil")
	}
	return &uploadDao{ctx}
}
