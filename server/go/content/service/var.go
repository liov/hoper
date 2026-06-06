package service

import "go.opentelemetry.io/otel"

var (
	momentSvc  = &MomentService{}
	actionSvc  = &ActionService{}
	contentSvc = &ContentService{}
	Tracer     = otel.Tracer("service")
)

func GetMomentService() *MomentService {
	if momentSvc != nil {
		return momentSvc
	}
	momentSvc = new(MomentService)
	return momentSvc
}

func GetActionService() *ActionService {
	if actionSvc != nil {
		return actionSvc
	}
	actionSvc = new(ActionService)
	return actionSvc
}

func GetContentService() *ContentService {
	if contentSvc != nil {
		return contentSvc
	}
	contentSvc = new(ContentService)
	return contentSvc
}
