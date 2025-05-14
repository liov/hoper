package data

import (
	"github.com/hopeio/context/httpctx"
	"log"
)

type chatDao struct {
	*httpctx.Context
}

func GetDao(ctx *httpctx.Context) *chatDao {
	if ctx == nil {
		log.Fatal("ctx can't nil")
	}
	return &chatDao{ctx}
}
