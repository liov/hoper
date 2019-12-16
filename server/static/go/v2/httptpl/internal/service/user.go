package service

import (
	"context"
	"time"

	"github.com/liov/hoper/go/v2/httptpl/internal/grpcclient"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/utils/log"
)

type UserService struct{}

func (*UserService) VerificationCode(req *model.VerifyReq) *model.VerifyRep {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	rep, err := grpcclient.UserClient.Verify(ctx, req)
	if err != nil {
		log.Errorf("could not greet: %v", err)
	}
	return rep
}

func (*UserService) Add(req *model.SignupReq) *model.SignupRep {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	rep, err := grpcclient.UserClient.Signup(ctx, req)
	if err != nil {
		log.Errorf("could not greet: %v", err)
	}
	return rep
}
