package service

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

import (
	"context"

	"github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/protobuf/utils/response"
)

type Resolver struct{}

func (r *mutationResolver) SignupVerify(ctx context.Context, in *user.SingUpVerifyReq) (*response.TinyRep, error) {
	panic("not implemented")
}

func (r *mutationResolver) Signup(ctx context.Context, in *user.SignupReq) (*user.SignupRep, error) {
	panic("not implemented")
}

func (r *mutationResolver) Active(ctx context.Context, in *user.ActiveReq) (*response.TinyRep, error) {
	panic("not implemented")
}

func (r *mutationResolver) Edit(ctx context.Context, in *user.EditReq) (*response.TinyRep, error) {
	panic("not implemented")
}

func (r *mutationResolver) ForgetPassword(ctx context.Context, in *user.LoginReq) (*response.TinyRep, error) {
	panic("not implemented")
}

func (r *mutationResolver) ResetPassword(ctx context.Context, in *user.ResetPasswordReq) (*response.TinyRep, error) {
	panic("not implemented")
}

func (r *queryResolver) VerifyCode(ctx context.Context) (*response.CommonRep, error) {
	panic("not implemented")
}

func (r *queryResolver) Login(ctx context.Context, in *user.LoginReq) (*user.LoginRep, error) {
	panic("not implemented")
}

func (r *queryResolver) Logout(ctx context.Context) (*user.LogoutRep, error) {
	panic("not implemented")
}

func (r *queryResolver) AuthInfo(ctx context.Context) (*user.UserMainInfo, error) {
	panic("not implemented")
}

func (r *queryResolver) GetUser(ctx context.Context, in *user.GetReq) (*user.GetRep, error) {
	panic("not implemented")
}

func (r *queryResolver) ActionLogList(ctx context.Context, in *user.ActionLogListReq) (*user.ActionLogListRep, error) {
	panic("not implemented")
}

// Mutation returns user.MutationResolver implementation.
func (r *Resolver) Mutation() user.MutationResolver { return &mutationResolver{r} }

// Query returns user.QueryResolver implementation.
func (r *Resolver) Query() user.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
