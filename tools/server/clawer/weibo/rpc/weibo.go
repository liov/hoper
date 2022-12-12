package rpc

type Weibo struct {
	CardlistInfo CardlistInfo `json:"cardlistInfo"`
	Cards        []*CardGroup `json:"cards"`
}

type Card struct {
	CardType     int          `json:"card_type"`
	ShowType     int          `json:"show_type"`
	CardStyle    int          `json:"card_style,omitempty"`
	Title        string       `json:"title"`
	CardGroup    []*CardGroup `json:"card_group,omitempty"`
	CardTypeName string       `json:"card_type_name,omitempty"`
	Itemid       string       `json:"itemid,omitempty"`
	Scheme       string       `json:"scheme,omitempty"`
	Mblog        *Mblog       `json:"mblog,omitempty"`
}

type Visible struct {
	Type      int    `json:"type"`
	ListId    int    `json:"list_id"`
	ListIdstr string `json:"list_idstr,omitempty"`
}

type Mblog struct {
	Visible                  Visible   `json:"visible"`
	CreatedAt                string    `json:"created_at"`
	Id                       string    `json:"id"`
	Mid                      string    `json:"mid"`
	CanEdit                  bool      `json:"can_edit"`
	ShowAdditionalIndication int       `json:"show_additional_indication"`
	Text                     string    `json:"text"`
	Source                   string    `json:"source"`
	Favorited                bool      `json:"favorited"`
	PicIds                   []string  `json:"pic_ids"`
	IsPaid                   bool      `json:"is_paid"`
	MblogVipType             int       `json:"mblog_vip_type"`
	User                     *UserInfo `json:"user"`
	PageInfo                 *PageInfo `json:"page_info,omitempty"`
	Pid                      int64     `json:"pid,omitempty"`
	Pidstr                   string    `json:"pidstr,omitempty"`
	RetweetedStatus          *Mblog    `json:"retweeted_status,omitempty"`
	RepostsCount             int       `json:"reposts_count"`
	CommentsCount            int       `json:"comments_count"`
	ReprintCmtCount          int       `json:"reprint_cmt_count"`
	AttitudesCount           int       `json:"attitudes_count"`
	PendingApprovalCount     int       `json:"pending_approval_count"`
	IsLongText               bool      `json:"isLongText"`
	Mlevel                   int       `json:"mlevel"`
	ShowMlevel               int       `json:"show_mlevel"`
	HotPage                  struct {
		Fid            string `json:"fid"`
		FeedDetailType int    `json:"feed_detail_type"`
	} `json:"hot_page"`
	Mblogtype             int    `json:"mblogtype"`
	Rid                   string `json:"rid"`
	ExternSafe            int    `json:"extern_safe"`
	NumberDisplayStrategy struct {
		ApplyScenarioFlag    int    `json:"apply_scenario_flag"`
		DisplayTextMinNumber int    `json:"display_text_min_number"`
		DisplayText          string `json:"display_text"`
	} `json:"number_display_strategy"`
	ContentAuth       int `json:"content_auth"`
	CommentManageInfo struct {
		CommentPermissionType int `json:"comment_permission_type"`
		ApprovalCommentType   int `json:"approval_comment_type"`
		CommentSortType       int `json:"comment_sort_type"`
	} `json:"comment_manage_info"`
	RepostType        int      `json:"repost_type,omitempty"`
	PicNum            int      `json:"pic_num"`
	Pics              []*Pic   `json:"pics"`
	NewCommentStyle   int      `json:"new_comment_style"`
	AbSwitcher        int      `json:"ab_switcher"`
	RegionName        string   `json:"region_name"`
	RegionOpt         int      `json:"region_opt"`
	MblogMenuNewStyle int      `json:"mblog_menu_new_style"`
	RawText           string   `json:"raw_text,omitempty"`
	Bid               string   `json:"bid"`
	TextLength        int      `json:"textLength,omitempty"`
	SafeTags          int64    `json:"safe_tags,omitempty"`
	Mark              string   `json:"mark,omitempty"`
	LivePhoto         []string `json:"live_photo,omitempty"`
}
type Pic struct {
	Pid      string   `json:"pid"`
	Url      string   `json:"url"`
	Size     string   `json:"size"`
	Large    PicLarge `json:"large"`
	VideoSrc string   `json:"videoSrc"`
	Type     string   `json:"type"`
}

type PicLarge struct {
	Size string `json:"size"`
	Url  string `json:"url"`
}

type CardlistInfo struct {
	CanShared         int                `json:"can_shared"`
	ShowStyle         int                `json:"show_style"`
	TitleTop          string             `json:"title_top"`
	PageType          string             `json:"page_type"`
	CardlistHeadCards []CardlistHeadCard `json:"cardlist_head_cards"`
	SinceId           int64              `json:"since_id"`
}

type CardlistHeadCard struct {
	HeadType    int       `json:"head_type"`
	ChannelList []Channel `json:"channel_list"`
}

type Channel struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Containerid string `json:"containerid"`
	DefaultAdd  int    `json:"default_add"`
	MustShow    int    `json:"must_show"`
	Apipath     string `json:"apipath"`
}

type CardGroup struct {
	CardType     int    `json:"card_type"`
	CardTypeName string `json:"card_type_name"`
	Itemid       string `json:"itemid"`
	Scheme       string `json:"scheme"`
	Mblog        *Mblog `json:"mblog"`
	ShowType     int    `json:"show_type"`
	Title        string `json:"title"`
}

type WeiboList struct {
	Cards        []*Card      `json:"cards"`
	CardlistInfo CardlistInfo `json:"cardlistInfo"`
	Scheme       string       `json:"scheme"`
	ShowAppTips  int          `json:"showAppTips"`
}

type PageInfo struct {
	Type       string `json:"type"`
	ObjectType int    `json:"object_type"`
	UrlOri     string `json:"url_ori,omitempty"`
	PagePic    struct {
		Pid         string `json:"pid,omitempty"`
		Source      int    `json:"source,omitempty"`
		IsSelfCover int    `json:"is_self_cover,omitempty"`
		Type        int    `json:"type,omitempty"`
		Url         string `json:"url"`
	} `json:"page_pic"`
	PageUrl          string    `json:"page_url"`
	ObjectId         string    `json:"object_id,omitempty"`
	PageTitle        string    `json:"page_title"`
	Title            string    `json:"title,omitempty"`
	Content1         string    `json:"content1"`
	Content2         string    `json:"content2,omitempty"`
	VideoOrientation string    `json:"video_orientation,omitempty"`
	PlayCount        string    `json:"play_count,omitempty"`
	MediaInfo        MediaInfo `json:"media_info,omitempty"`
	Urls             VideoUrl  `json:"urls,omitempty"`
}

type MediaInfo struct {
	StreamUrl   string  `json:"stream_url"`
	StreamUrlHd string  `json:"stream_url_hd"`
	Duration    float64 `json:"duration"`
}

type VideoUrl struct {
	Mp4720PMp4 string `json:"mp4_720p_mp4"`
	Mp4HdMp4   string `json:"mp4_hd_mp4"`
	Mp4LdMp4   string `json:"mp4_ld_mp4"`
}

var badges = []string{
	"anniversary",
	"gongyi",
	"bind_taobao",
	"unread_pool",
	"unread_pool_ext",
	"dzwbqlx_2016",
	"cz_wed_2017",
	"panda",
	"user_name_certificate",
	"suishoupai_2018",
	"super_star_2018",
	"dailv_2018",
	"asiad_2018",
	"qixi_2018",
	"lol_s8",
	"hongbaofei_2019",
	"fu_2019",
	"suishoupai_2019",
	"china_2019",
	"hongkong_2019",
	"hongbao_2020",
	"feiyan_2020",
	"hongbaofeifuniu_2021",
	"hongbaofeijika_2021",
	"biyeji_2021",
	"aoyun_2021",
	"dailu_2021",
	"kaixue21_2021",
	"newdongaohui_2022",
	"biyeji_2022",
	"gaokao_2022",
}
