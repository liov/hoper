package rpc

type User struct {
	IsVideoCoverStyle int      `json:"isVideoCoverStyle"`
	IsStarStyle       int      `json:"isStarStyle"`
	UserInfo          UserInfo `json:"userInfo"`
	FansScheme        string   `json:"fans_scheme"`
	FollowScheme      string   `json:"follow_scheme"`
	TabsInfo          struct {
		SelectedTab int `json:"selectedTab"`
		Tabs        []struct {
			Id          int    `json:"id"`
			TabKey      string `json:"tabKey"`
			MustShow    int    `json:"must_show"`
			Hidden      int    `json:"hidden"`
			Title       string `json:"title"`
			TabType     string `json:"tab_type"`
			Containerid string `json:"containerid"`
			Apipath     string `json:"apipath,omitempty"`
			TabIcon     string `json:"tab_icon,omitempty"`
			TabIconDark string `json:"tab_icon_dark,omitempty"`
			Url         string `json:"url,omitempty"`
		} `json:"tabs"`
	} `json:"tabsInfo"`
	ProfileExt  string `json:"profile_ext"`
	Scheme      string `json:"scheme"`
	ShowAppTips int    `json:"showAppTips"`
}

type UserInfo struct {
	Id                int            `json:"id"`
	ScreenName        string         `json:"screen_name"`
	ProfileImageUrl   string         `json:"profile_image_url"`
	ProfileUrl        string         `json:"profile_url"`
	StatusesCount     int            `json:"statuses_count"`
	Verified          bool           `json:"verified"`
	VerifiedType      int            `json:"verified_type"`
	CloseBlueV        bool           `json:"close_blue_v"`
	Description       string         `json:"description"`
	Gender            string         `json:"gender"`
	Mbtype            int            `json:"mbtype"`
	Svip              int            `json:"svip"`
	Urank             int            `json:"urank"`
	Mbrank            int            `json:"mbrank"`
	FollowMe          bool           `json:"follow_me"`
	Following         bool           `json:"following"`
	FollowCount       int            `json:"follow_count"`
	FollowersCount    string         `json:"followers_count"`
	FollowersCountStr string         `json:"followers_count_str"`
	CoverImagePhone   string         `json:"cover_image_phone"`
	AvatarHd          string         `json:"avatar_hd"`
	ToolbarMenus      []*ToolbarMenu `json:"toolbar_menus"`
}

type ToolbarMenu struct {
	Type     string           `json:"type"`
	Name     string           `json:"name"`
	Pic      string           `json:"pic"`
	Params   ToolbarMenuParam `json:"params"`
	Scheme   string           `json:"scheme,omitempty"`
	UserInfo *SubUserInfo     `json:"userInfo,omitempty"`
}

type ToolbarMenuParam struct {
	Scheme    string `json:"scheme,omitempty"`
	Uid       int64  `json:"uid,omitempty"`
	Extparams struct {
		Followcardid string `json:"followcardid"`
	} `json:"extparams,omitempty"`
}

type SubUserInfo struct {
	Id                  int64  `json:"id"`
	Idstr               string `json:"idstr"`
	ScreenName          string `json:"screen_name"`
	ProfileImageUrl     string `json:"profile_image_url"`
	Following           bool   `json:"following"`
	Verified            bool   `json:"verified"`
	VerifiedType        int    `json:"verified_type"`
	Remark              string `json:"remark"`
	AvatarLarge         string `json:"avatar_large"`
	AvatarHd            string `json:"avatar_hd"`
	FollowMe            bool   `json:"follow_me"`
	Mbtype              int    `json:"mbtype"`
	Mbrank              int    `json:"mbrank"`
	Level               int    `json:"level"`
	Type                int    `json:"type"`
	StoryReadState      int    `json:"story_read_state"`
	AllowMsg            int    `json:"allow_msg"`
	FriendshipsRelation int    `json:"friendships_relation"`
	CloseFriendsType    int    `json:"close_friends_type"`
	SpecialFollow       bool   `json:"special_follow"`
}
