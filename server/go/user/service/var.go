package service

import "go.opentelemetry.io/otel"

var (
	userSvc  = &UserService{}
	oauthSvc *OauthService

	Tracer = otel.Tracer("service")
	Meter  = otel.Meter("service")
)

func GetUserService() *UserService {
	if userSvc != nil {
		return userSvc
	}
	userSvc = new(UserService)
	return userSvc
}
