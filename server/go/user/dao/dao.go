package dao

import (
	"github.com/hopeio/lemon/context/http_context"
	"log"
)

type userDao struct {
	*http_context.Context
}

func GetDao(ctx *http_context.Context) *userDao {
	if ctx == nil {
		log.Fatal("ctx can't nil")
	}
	return &userDao{ctx}
}
