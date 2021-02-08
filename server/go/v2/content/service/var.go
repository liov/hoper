package service

var (
	momentSvc = &MomentService{}
)

func GetMomentService() *MomentService {
	if momentSvc != nil {
		return momentSvc
	}
	momentSvc = new(MomentService)
	return momentSvc
}