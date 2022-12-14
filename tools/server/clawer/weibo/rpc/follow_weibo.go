package rpc

type WeiboFollowsList struct {
	Statuses          []*Mblog `json:"statuses"`
	Hasvisible        bool     `json:"hasvisible"`
	PreviousCursor    int      `json:"previous_cursor"`
	NextCursor        int64    `json:"next_cursor"`
	PreviousCursorStr string   `json:"previous_cursor_str"`
	NextCursorStr     string   `json:"next_cursor_str"`
	TotalNumber       int      `json:"total_number"`
	Interval          int      `json:"interval"`
	UveBlank          int      `json:"uve_blank"`
	SinceId           int64    `json:"since_id"`
	SinceIdStr        string   `json:"since_id_str"`
	MaxId             int64    `json:"max_id"`
	MaxIdStr          string   `json:"max_id_str"`
	HasUnread         int      `json:"has_unread"`
}
