package service

import (
	"go.opentelemetry.io/otel"
)

const ScopeName = "github.com/liov/hoper/server/go/user"
var (
	userSvc  *UserService
	oauthSvc *OauthService

	Tracer = otel.Tracer(ScopeName)
	Meter  = otel.Meter(ScopeName)
)

func GetUserService() *UserService {
	if userSvc != nil {
		return userSvc
	}
	userSvc = new(UserService)
	return userSvc
}
