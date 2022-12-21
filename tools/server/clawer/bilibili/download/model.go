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
	PubAt     time.Time
	CreatedAt time.Time
}

func NewVideo(uid int, title string, aid, cid, page int, part string, pubdate int) *Video {
	return &Video{
		Uid:   uid,
		Title: title,
		Aid:   aid,
		Cid:   cid,
		Page:  page,
		Part:  part,
		PubAt: time.Unix(int64(pubdate), 0),
	}
}
