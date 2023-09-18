package data

import (
	"github.com/hopeio/lemon/context/http_context"
	"log"
)

type uploadDao struct {
	*http_context.Context
}

func GetDao(ctx *http_context.Context) *uploadDao {
	if ctx == nil {
		log.Fatal("ctx can't nil")
	}
	return &uploadDao{ctx}
}
