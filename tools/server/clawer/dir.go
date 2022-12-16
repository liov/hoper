package clawer

type Dir struct {
	Type     int    `json:"type"`
	UserId   int    `json:"userId"`
	UserName string `json:"userName"`
	FilePath string `json:"filePath"`
}
