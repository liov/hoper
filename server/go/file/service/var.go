package service

var (
	fileSvc = &FileService{}
)

func GetFileService() *FileService {
	if fileSvc != nil {
		return fileSvc
	}
	fileSvc = new(FileService)
	return fileSvc
}
