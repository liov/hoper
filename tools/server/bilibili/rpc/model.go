package rpc

type DescV2 struct {
	RawText string `json:"raw_text"`
	Type    int    `json:"type"`
	BizId   int    `json:"biz_id"`
}

type Rights struct {
	Bp            int `json:"bp"`
	Elec          int `json:"elec"`
	Download      int `json:"download"`
	Movie         int `json:"movie"`
	Pay           int `json:"pay"`
	Hd5           int `json:"hd5"`
	NoReprint     int `json:"no_reprint"`
	Autoplay      int `json:"autoplay"`
	UgcPay        int `json:"ugc_pay"`
	IsCooperation int `json:"is_cooperation"`
	UgcPayPreview int `json:"ugc_pay_preview"`
	NoBackground  int `json:"no_background"`
	CleanMode     int `json:"clean_mode"`
	IsSteinGate   int `json:"is_stein_gate"`
	Is360         int `json:"is_360"`
	NoShare       int `json:"no_share"`
	ArcPay        int `json:"arc_pay"`
	FreeWatch     int `json:"free_watch"`
}

type Stat struct {
	Aid        int    `json:"aid" gorm:"primaryKey"`
	View       int    `json:"view"`
	Danmaku    int    `json:"danmaku"`
	Reply      int    `json:"reply"`
	Favorite   int    `json:"favorite"`
	Coin       int    `json:"coin"`
	Share      int    `json:"share"`
	NowRank    int    `json:"now_rank"`
	HisRank    int    `json:"his_rank"`
	Like       int    `json:"like"`
	Dislike    int    `json:"dislike"`
	Evaluation string `json:"evaluation"`
	ArgueMsg   string `json:"argue_msg"`
}

type Dimension struct {
	Width  int `json:"width"`
	Height int `json:"height"`
	Rotate int `json:"rotate"`
}

type Page struct {
	Cid        int        `json:"cid" gorm:"primaryKey"`
	Page       int        `json:"page"`
	From       string     `json:"from"`
	Part       string     `json:"part"`
	Duration   int        `json:"duration"`
	Vid        string     `json:"vid"`
	WebLink    string     `json:"weblink"`
	Dimension  *Dimension `json:"dimension"`
	FirstFrame string     `json:"first_frame"`
}

type Owner struct {
	Mid  int    `json:"mid" gorm:"primaryKey"`
	Name string `json:"name"`
	Face string `json:"face"`
}

type ViewInfo struct {
	Bvid               string      `json:"bvid" gorm:"index:idx_bvid,unique"`
	Aid                int         `json:"aid" gorm:"primaryKey"`
	Videos             int         `json:"videos"`
	Tid                int         `json:"tid"`
	Tname              string      `json:"tname"`
	Copyright          int         `json:"copyright"`
	Pic                string      `json:"pic"`
	Title              string      `json:"title"`
	PubDate            int         `json:"pubdate"`
	Ctime              int         `json:"ctime"`
	Desc               string      `json:"desc"`
	DescV2             []*DescV2   `json:"desc_v2" gorm:"-"`
	State              int         `json:"state"`
	Duration           int         `json:"duration"`
	Rights             *Rights     `json:"rights"`
	Owner              *Owner      `json:"owner" gorm:"-"`
	OwnerMid           int         `json:"-" gorm:"index"`
	Stat               *Stat       `json:"stat"`
	Dynamic            string      `json:"dynamic"`
	Cid                int         `json:"cid"`
	Dimension          *Dimension  `json:"dimension"`
	Premiere           interface{} `json:"premiere,omitempty"`
	TeenageMode        int         `json:"teenage_mode"`
	IsChargeableSeason bool        `json:"is_chargeable_season"`
	IsStory            bool        `json:"is_story"`
	NoCache            bool        `json:"no_cache"`
	Pages              []*Page     `json:"pages"`
	Subtitle           struct {
		AllowSubmit bool          `json:"allow_submit"`
		List        []interface{} `json:"list"`
	} `json:"subtitle" gorm:"-"`
	IsSeasonDisplay bool `json:"is_season_display"`
	UserGarb        struct {
		UrlImageAniCut string `json:"url_image_ani_cut"`
	} `json:"user_garb" gorm:"-"`
	HonorReply struct {
	} `json:"honor_reply,omitempty" gorm:"-"`
}

func (*ViewInfo) TableName() string {
	return "view"
}

type NavInfo struct {
	IsLogin            bool           `json:"isLogin"`
	EmailVerified      int            `json:"email_verified"`
	Face               string         `json:"face"`
	FaceNft            int            `json:"face_nft"`
	FaceNftType        int            `json:"face_nft_type"`
	LevelInfo          *LevelInfo     `json:"level_info"`
	Mid                int            `json:"mid"`
	MobileVerified     int            `json:"mobile_verified"`
	Money              int            `json:"money"`
	Moral              int            `json:"moral"`
	Official           *Official      `json:"official"`
	OfficialVerify     OfficialVerify `json:"officialVerify"`
	Pendant            *Pendant       `json:"pendant"`
	Scores             int            `json:"scores"`
	Uname              string         `json:"uname"`
	VipDueDate         int64          `json:"vipDueDate"`
	VipStatus          int            `json:"vipStatus"`
	VipType            int            `json:"vipType"`
	VipPayType         int            `json:"vip_pay_type"`
	VipThemeType       int            `json:"vip_theme_type"`
	VipLabel           *Label         `json:"vip_label"`
	VipAvatarSubscript int            `json:"vip_avatar_subscript"`
	VipNicknameColor   string         `json:"vip_nickname_color"`
	Vip                *Vip           `json:"vip"`
	Wallet             *Wallet        `json:"wallet"`
	HasShop            bool           `json:"has_shop"`
	ShopUrl            string         `json:"shop_url"`
	AllowanceCount     int            `json:"allowance_count"`
	AnswerStatus       int            `json:"answer_status"`
	IsSeniorMember     int            `json:"is_senior_member"`
}

type OfficialVerify struct {
	Type int    `json:"type"`
	Desc string `json:"desc"`
}

type LevelInfo struct {
	CurrentLevel int `json:"current_level"`
	CurrentMin   int `json:"current_min"`
	CurrentExp   int `json:"current_exp"`
	NextExp      int `json:"next_exp"`
}

type Official struct {
	Role  int    `json:"role"`
	Title string `json:"title"`
	Desc  string `json:"desc"`
	Type  int    `json:"type"`
}

type Pendant struct {
	Pid               int    `json:"pid"`
	Name              string `json:"name"`
	Image             string `json:"image"`
	Expire            int    `json:"expire"`
	ImageEnhance      string `json:"image_enhance"`
	ImageEnhanceFrame string `json:"image_enhance_frame"`
}

type Wallet struct {
	Mid           int `json:"mid"`
	BcoinBalance  int `json:"bcoin_balance"`
	CouponBalance int `json:"coupon_balance"`
	CouponDueTime int `json:"coupon_due_time"`
}

type Vip struct {
	Type               int    `json:"type"`
	Status             int    `json:"status"`
	DueDate            int64  `json:"due_date"`
	VipPayType         int    `json:"vip_pay_type"`
	ThemeType          int    `json:"theme_type"`
	Label              *Label `json:"label"`
	AvatarSubscript    int    `json:"avatar_subscript"`
	NicknameColor      string `json:"nickname_color"`
	Role               int    `json:"role"`
	AvatarSubscriptUrl string `json:"avatar_subscript_url"`
	TvVipStatus        int    `json:"tv_vip_status"`
	TvVipPayType       int    `json:"tv_vip_pay_type"`
}

type Label struct {
	Path                  string `json:"path"`
	Text                  string `json:"text"`
	LabelTheme            string `json:"label_theme"`
	TextColor             string `json:"text_color"`
	BgStyle               int    `json:"bg_style"`
	BgColor               string `json:"bg_color"`
	BorderColor           string `json:"border_color"`
	UseImgLabel           bool   `json:"use_img_label"`
	ImgLabelUriHans       string `json:"img_label_uri_hans"`
	ImgLabelUriHant       string `json:"img_label_uri_hant"`
	ImgLabelUriHansStatic string `json:"img_label_uri_hans_static"`
	ImgLabelUriHantStatic string `json:"img_label_uri_hant_static"`
}

type VideoInfo struct {
	Aid               int      `json:"-" gorm:"index"`
	Cid               int      `json:"-" gorm:"index"`
	From              string   `json:"from,omitempty"`
	Result            string   `json:"result,omitempty"`
	Quality           int      `json:"quality,omitempty"`
	Format            string   `json:"format,omitempty"`
	Timelength        int      `json:"timelength"`
	AcceptFormat      string   `json:"accept_format"`
	AcceptDescription []string `json:"accept_description,omitempty"`
	AcceptQuality     []int    `json:"accept_quality"`
	VideoCodecid      int      `json:"video_codecid" gorm:"primaryKey"`
	VideoProject      bool     `json:"video_project"`
	SeekParam         string   `json:"seek_param,omitempty"`
	SeekType          string   `json:"seek_type,omitempty"`
	Durl              []*Durl  `json:"durl" gorm:"-"`
}

func (v *VideoInfo) JsonClean() {
	v.From = ""
	v.Result = ""
	v.Quality = 0
	v.Format = ""
	v.SeekParam = ""
	v.SeekType = ""
	v.AcceptDescription = nil
	for _, durl := range v.Durl {
		durl.Url = ""
		durl.BackupUrl = nil
	}
}

type Durl struct {
	Order     int      `json:"order"`
	Length    int      `json:"length"`
	Size      int      `json:"size"`
	Url       string   `json:"url,omitempty"`
	BackupUrl []string `json:"backup_url,omitempty"`
}

type FavList struct {
	Info    FavInfo  `json:"info"`
	Medias  []*Media `json:"medias"`
	HasMore bool     `json:"has_more"`
}

type FavInfo struct {
	Id    int    `json:"id"`
	Fid   int    `json:"fid"`
	Mid   int    `json:"mid"`
	Attr  int    `json:"attr"`
	Title string `json:"title"`
	Cover string `json:"cover"`
	Upper struct {
		Mid       int    `json:"mid"`
		Name      string `json:"name"`
		Face      string `json:"face"`
		Followed  bool   `json:"followed"`
		VipType   int    `json:"vip_type"`
		VipStatue int    `json:"vip_statue"`
	} `json:"upper"`
	CoverType int `json:"cover_type"`
	CntInfo   struct {
		Collect int `json:"collect"`
		Play    int `json:"play"`
		ThumbUp int `json:"thumb_up"`
		Share   int `json:"share"`
	} `json:"cnt_info"`
	Type       int    `json:"type"`
	Intro      string `json:"intro"`
	Ctime      int    `json:"ctime"`
	Mtime      int    `json:"mtime"`
	State      int    `json:"state"`
	FavState   int    `json:"fav_state"`
	LikeState  int    `json:"like_state"`
	MediaCount int    `json:"media_count"`
}

type Media struct {
	Id       int    `json:"id"`
	Type     int    `json:"type"`
	Title    string `json:"title"`
	Cover    string `json:"cover"`
	Intro    string `json:"intro"`
	Page     int    `json:"page"`
	Duration int    `json:"duration"`
	Upper    struct {
		Mid  int    `json:"mid"`
		Name string `json:"name"`
		Face string `json:"face"`
	} `json:"upper"`
	Attr    int `json:"attr"`
	CntInfo struct {
		Collect int `json:"collect"`
		Play    int `json:"play"`
		Danmaku int `json:"danmaku"`
	} `json:"cnt_info"`
	Link    string      `json:"link"`
	Ctime   int         `json:"ctime"`
	Pubtime int         `json:"pubtime"`
	FavTime int         `json:"fav_time"`
	BvId    string      `json:"bv_id"`
	Bvid    string      `json:"bvid"`
	Season  interface{} `json:"season"`
	Ogv     interface{} `json:"ogv"`
	Ugc     struct {
		FirstCid int `json:"first_cid"`
	} `json:"ugc"`
}

type UpSpaceList struct {
	List List `json:"list"`
	Page struct {
		Pn    int `json:"pn"`
		Ps    int `json:"ps"`
		Count int `json:"count"`
	} `json:"page"`
	EpisodicButton struct {
		Text string `json:"text"`
		Uri  string `json:"uri"`
	} `json:"episodic_button"`
}

type List struct {
	//Tlist map[string]Tag `json:"tlist"`
	Vlist []*Video `json:"vlist"`
}

type Tag struct {
	Tid   int    `json:"tid"`
	Count int    `json:"count"`
	Name  string `json:"name"`
}

type Video struct {
	Comment        int    `json:"comment"`
	Typeid         int    `json:"typeid"`
	Play           int    `json:"play"`
	Pic            string `json:"pic"`
	Subtitle       string `json:"subtitle"`
	Description    string `json:"description"`
	Copyright      string `json:"copyright"`
	Title          string `json:"title"`
	Review         int    `json:"review"`
	Author         string `json:"author"`
	Mid            int    `json:"mid"`
	Created        int    `json:"created"`
	Length         string `json:"length"`
	VideoReview    int    `json:"video_review"`
	Aid            int    `json:"aid"`
	Bvid           string `json:"bvid"`
	HideClick      bool   `json:"hide_click"`
	IsPay          int    `json:"is_pay"`
	IsUnionVideo   int    `json:"is_union_video"`
	IsSteinsGate   int    `json:"is_steins_gate"`
	IsLivePlayback int    `json:"is_live_playback"`
}
