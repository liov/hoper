package dao

type User struct {
	Id               string `json:"id"`
	ScreenName       string `json:"screen_name"`
	Gender           string `json:"gender"`
	Birthday         string `json:"birthday"`
	Location         string `json:"location"`
	Education        string `json:"education"`
	Company          string `json:"company"`
	RegistrationTime string `json:"registration_time"`
	Sunshine         string `json:"sunshine"`
	StatusesCount    int    `json:"statuses_count"`
	FollowersCount   int    `json:"followers_count"`
	FollowCount      int    `json:"follow_count"`
	Description      string `json:"description"`
	ProfileUrl       string `json:"profile_url"`
	ProfileImageUrl  string `json:"profile_image_url"`
	AvatarHd         string `json:"avatar_hd"`
	Urank            int    `json:"urank"`
	Mbrank           int    `json:"mbrank"`
	Verified         bool   `json:"verified"`
	VerifiedType     int    `json:"verified_type"`
	VerifiedReason   string `json:"verified_reason"`
}

type Retweet struct {
	UserId         int    `json:"user_id"`
	ScreenName     string `json:"screen_name"`
	Id             int64  `json:"id"`
	Bid            string `json:"bid"`
	Text           string `json:"text"`
	Pics           string `json:"pics"`
	VideoUrl       string `json:"video_url"`
	Location       string `json:"location"`
	CreatedAt      string `json:"created_at"`
	Source         string `json:"source"`
	AttitudesCount int    `json:"attitudes_count"`
	CommentsCount  int    `json:"comments_count"`
	RepostsCount   int    `json:"reposts_count"`
	Topics         string `json:"topics"`
	AtUsers        string `json:"at_users"`
}

type Weibo struct {
}
