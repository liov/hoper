package service

import (
	"reflect"
	"testing"

	"github.com/liov/hoper/go/v2/utils/log"
)

func TestUserService_Signup(t *testing.T) {
	user := UserService{}
	log.Info(reflect.TypeOf(user).Size())
}
