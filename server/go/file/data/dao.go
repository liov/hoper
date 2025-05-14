package data

import (
	"github.com/hopeio/context/httpctx"
	"log"
)

type uploadDao struct {
	*httpctx.Context
}

func GetDao(ctx *httpctx.Context) *uploadDao {
	if ctx == nil {
		log.Fatal("ctx can't nil")
	}
	return &uploadDao{ctx}
}
