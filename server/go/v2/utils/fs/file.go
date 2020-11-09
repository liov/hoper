package fs

type File struct {
	ID           uint64 `gorm:"primary_key" json:"id"`
	FileName     string `gorm:"type:varchar(100);not null" json:"file_name"`
	OriginalName string `gorm:"type:varchar(100);not null" json:"original_name"`
	URL          string `json:"url"`
	MD5          string `gorm:"type:varchar(32)" json:"md5"`
	Mime         string `json:"mime"`
	Size         uint64 `json:"size"`
}
