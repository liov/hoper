package dao

import "time"

const (
	TableNameWeibo   = "weibo.weibo"
	TableNameUser    = "weibo.user"
	TableNameComment = "weibo.comment"
)

type User struct {
	Id          int    `json:"id"`
	ScreenName  string `json:"screen_name"`
	Gender      string `json:"gender"`
	Sunshine    string `json:"sunshine"`
	Description string `json:"description"`
	AvatarHd    string `json:"avatar_hd"`
}

func (w *User) TableName() string {
	return TableNameUser
}

type Weibo struct {
	Id          int    `json:"id"`
	BId         string `json:"bid"`
	UserId      int    `json:"user_id"`
	Text        string
	CreatedAt   time.Time `json:"created_at"`
	Source      string    `json:"source"`
	Pics        string
	Video       string
	RetweetedId int
}

func (w *Weibo) TableName() string {
	return TableNameWeibo
}

type Comment struct {
	Id        int       `json:"id"`
	UserId    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	Source    string    `json:"source"`
	Text      string    `json:"text"`
	ReplyId   int64     `json:"reply_id,omitempty"`
	ReplyText string    `json:"reply_text,omitempty"`
}

func (w *Comment) TableName() string {
	return TableNameComment
}
