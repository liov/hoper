package model

type Video struct {
	Aid               int      `json:"aid" gorm:"primaryKey"`
	Cid               int      `json:"cid" gorm:"primaryKey"`
	Timelength        int      `json:"timelength"`
	AcceptFormat      string   `json:"accept_format"`
	AcceptDescription []string `json:"accept_description"`
	AcceptQuality     []int    `json:"accept_quality"`
	VideoCodecid      int      `json:"video_codecid" gorm:"primaryKey"`
	VideoProject      bool     `json:"video_project"`
	SeekParam         string   `json:"seek_param"`
	SeekType          string   `json:"seek_type"`
}
