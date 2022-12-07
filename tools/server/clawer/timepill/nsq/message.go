package nsq

import "tools/clawer/timepill/model"

type Diary struct {
	UserId   int    `json:"user_id"`
	PhotoUrl string `json:"photo_url"`
	Created  string `json:"created"`
}

type Cover struct {
	Type model.CoverType `json:"type"`
	Url  string          `json:"url"`
}
