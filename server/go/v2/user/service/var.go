package service

var (
	userSvc  = &UserService{}
	oauthSvc *OauthService
)

func GetUserService() *UserService {
	if userSvc != nil {
		return userSvc
	}
	userSvc = new(UserService)
	return userSvc
}