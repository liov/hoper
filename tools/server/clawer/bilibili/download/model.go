package download

import "time"

type Video struct {
	Uid       int
	Title     string
	Aid       int
	Cid       int
	Page      int
	Part      string
	Quality   int
	Record    int
	CodecId   int
	CreatedAt time.Time
}

func NewVideo(uid int, title string, aid, cid, page int, part string) *Video {
	return &Video{
		Uid:   uid,
		Title: title,
		Aid:   aid,
		Cid:   cid,
		Page:  page,
		Part:  part,
	}
}
