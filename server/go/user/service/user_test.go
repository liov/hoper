package service

import (
	"reflect"
	"testing"

	"github.com/hopeio/dora/utils/log"
)

func TestUserService_Signup(t *testing.T) {
	user := UserService{}
	log.Info(reflect.TypeOf(user).Size())
}
