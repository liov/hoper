package service

import "go.opentelemetry.io/otel"

var (
	fileSvc = &FileService{}
	Tracer  = otel.Tracer("file")
)

func GetFileService() *FileService {
	if fileSvc != nil {
		return fileSvc
	}
	fileSvc = new(FileService)
	return fileSvc
}
