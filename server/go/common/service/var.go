package service

var (
	commonSvc = &CommonService{}
)

func GetMomentService() *CommonService {
	if commonSvc != nil {
		return commonSvc
	}
	commonSvc = new(CommonService)
	return commonSvc
}
