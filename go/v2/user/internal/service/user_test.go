package service

import (
	"reflect"
	"testing"

	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/utils/log"
)

func TestUserService_Signup(t *testing.T) {
	user:=model.User{}
	log.Info(reflect.TypeOf(user).Size())
}
