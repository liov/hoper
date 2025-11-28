package service

var (
	commonSvc = &CommonService{}
)

func GetCommonService() *CommonService {
	if commonSvc != nil {
		return commonSvc
	}
	commonSvc = new(CommonService)
	return commonSvc
}
