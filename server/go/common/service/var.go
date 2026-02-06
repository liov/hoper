package service

import "go.opentelemetry.io/otel"

var (
	commonSvc = &CommonService{}
	Trancer   = otel.Tracer("service")
)

func GetCommonService() *CommonService {
	if commonSvc != nil {
		return commonSvc
	}
	commonSvc = new(CommonService)
	return commonSvc
}
