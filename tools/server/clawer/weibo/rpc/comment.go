package rpc

type CommentList struct {
	Data        []*Comment `json:"data"`
	TotalNumber int        `json:"total_number"`
	Max         int        `json:"max"`
}

type Comment struct {
	Id         int       `json:"id"`
	CreatedAt  string    `json:"created_at"`
	Source     string    `json:"source"`
	User       *UserInfo `json:"user"`
	Text       string    `json:"text"`
	LikeCounts int       `json:"like_counts"`
	Liked      bool      `json:"liked"`
	ReplyId    int64     `json:"reply_id,omitempty"`
	ReplyText  string    `json:"reply_text,omitempty"`
}
