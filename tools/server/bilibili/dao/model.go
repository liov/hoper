package dao

const (
	Schema         = "bilibili."
	TableNameView  = Schema + "view"
	TableNameVideo = Schema + "video"
)

type View struct {
	Bvid        string `json:"bvid" gorm:"index:idx_bvid,unique"`
	Aid         int    `json:"aid" gorm:"primaryKey"`
	Data        []byte `json:"data" gorm:"type:jsonb"`
	CoverRecord bool
}

func (v *View) TableName() string {
	return TableNameView
}

type Video struct {
	Aid    int    `json:"aid" gorm:"primaryKey"`
	Cid    int    `json:"cid" gorm:"primaryKey"`
	Data   []byte `json:"data" gorm:"type:jsonb"`
	Record bool
}

func (v *Video) TableName() string {
	return TableNameVideo
}
