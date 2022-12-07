package rpc

type PicCards struct {
	CardlistInfo PicCardlistInfo `json:"cardlistInfo"`
	Cards        []*PicCardGroup `json:"cards"`
	ShowAppTips  int             `json:"showAppTips"`
}

type PicCard struct {
	CardType     int             `json:"card_type"`
	CardTypeName string          `json:"card_type_name"`
	Title        string          `json:"title"`
	Itemid       string          `json:"itemid"`
	CardGroup    []*PicCardGroup `json:"card_group"`
	Buttontitle  string          `json:"buttontitle,omitempty"`
	DisplayArrow int             `json:"display_arrow"`
}

type PicCardlistInfo struct {
	VP          int           `json:"v_p"`
	Containerid string        `json:"containerid"`
	TitleTop    string        `json:"title_top"`
	ShowStyle   int           `json:"show_style"`
	Total       int           `json:"total"`
	FilterGroup []interface{} `json:"filter_group"`
	Page        interface{}   `json:"page"`
}

type PicCardGroup struct {
	CardType     int     `json:"card_type"`
	CardTypeName string  `json:"card_type_name"`
	Itemid       string  `json:"itemid"`
	Pics         []*Pics `json:"pics,omitempty"`
	Pic          string  `json:"pic,omitempty"`
	TitleSub     string  `json:"title_sub,omitempty"`
	Desc1        string  `json:"desc1,omitempty"`
	Desc2        string  `json:"desc2,omitempty"`
	DisplayArrow int     `json:"display_arrow,omitempty"`
}

type Pics struct {
	PicSmall    string    `json:"pic_small"`
	PicMiddle   string    `json:"pic_middle"`
	PicBig      string    `json:"pic_big"`
	PicMw2000   string    `json:"pic_mw2000"`
	Mblog       TinyMblog `json:"mblog"`
	ObjectId    string    `json:"object_id"`
	Actionlog   string    `json:"actionlog"`
	PhotoTag    int       `json:"photo_tag"`
	PicId       string    `json:"pic_id"`
	Savedisable int       `json:"savedisable"`
	Pic         string    `json:"pic"`
	Video       string    `json:"video,omitempty"`
	Type        string    `json:"type,omitempty"`
}

type TinyMblog struct {
	Id           string              `json:"id"`
	Mid          string              `json:"mid"`
	Text         string              `json:"text"`
	IsLongText   bool                `json:"isLongText"`
	Scheme       string              `json:"scheme"`
	PicIdsX      interface{}         `json:"pic_ids_x"`
	PicIds       []string            `json:"pic_ids"`
	PicInfos     map[string]PicInfos `json:"pic_infos"`
	MblogVipType int                 `json:"mblog_vip_type"`
	IsPaid       bool                `json:"is_paid"`
}

type PicInfos struct {
	Thumbnail PicInfo `json:"thumbnail"`
	Bmiddle   PicInfo `json:"bmiddle"`
	Large     PicInfo `json:"large"`
	Original  PicInfo `json:"original"`
	Mw2000    PicInfo `json:"mw2000"`
	ObjectId  string  `json:"object_id"`
	PicId     string  `json:"pic_id"`
	PhotoTag  int     `json:"photo_tag"`
	Type      string  `json:"type"`
	PicStatus int     `json:"pic_status"`
}

type PicInfo struct {
	Url string `json:"url"`
}
