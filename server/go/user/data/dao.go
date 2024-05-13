package data

import (
	"github.com/hopeio/cherry/context/httpctx"
	"log"
)

type userDao struct {
	*httpctx.Context
}

func GetDao(ctx *httpctx.Context) *userDao {
	if ctx == nil {
		log.Fatal("ctx can't nil")
	}
	return &userDao{ctx}
}
