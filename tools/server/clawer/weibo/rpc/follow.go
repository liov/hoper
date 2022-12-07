package rpc

type Follow struct {
	CardlistInfo struct {
		Containerid string `json:"containerid"`
		TitleTop    string `json:"title_top"`
		ShowStyle   string `json:"show_style"`
		Total       int    `json:"total"`
		Page        int    `json:"page"`
	} `json:"cardlistInfo"`
	Cards  []*FollowCard `json:"cards"`
	Scheme string        `json:"scheme"`
}

type FollowCard struct {
	CardStyle int                `json:"card_style,omitempty"`
	CardType  int                `json:"card_type"`
	Itemid    string             `json:"itemid"`
	CardGroup []*FollowCardGroup `json:"card_group"`
	Title     string             `json:"title,omitempty"`
}

type FollowCardGroup struct {
	CardType        int         `json:"card_type"`
	Scheme          string      `json:"scheme"`
	DisplayArrow    int         `json:"display_arrow,omitempty"`
	TitleExtraText  string      `json:"title_extra_text,omitempty"`
	Itemid          string      `json:"itemid"`
	Desc            string      `json:"desc,omitempty"`
	Actionlog       Actionlog   `json:"actionlog,omitempty"`
	CardTypeName    string      `json:"card_type_name,omitempty"`
	Title           string      `json:"title,omitempty"`
	WeiboNeed       string      `json:"weibo_need,omitempty"`
	Elements        []*Element  `json:"elements,omitempty"`
	Users           []*UserInfo `json:"users,omitempty"`
	ShowType        int         `json:"show_type,omitempty"`
	BackgroundColor int         `json:"background_color,omitempty"`
	Openurl         string      `json:"openurl,omitempty"`
	RecomRemark     string      `json:"recom_remark,omitempty"`
	Recommend       string      `json:"recommend,omitempty"`
	Desc1           string      `json:"desc1,omitempty"`
	Desc2           string      `json:"desc2,omitempty"`
	User            UserInfo    `json:"user,omitempty"`
	Buttons         []struct {
		Type       string `json:"type"`
		SubType    int    `json:"sub_type"`
		Name       string `json:"name"`
		SkipFormat int    `json:"skip_format"`
		Params     struct {
			Uid        int64 `json:"uid"`
			NeedFollow int   `json:"need_follow"`
			TrendExt   int64 `json:"trend_ext"`
			TrendType  int   `json:"trend_type"`
			Itemid     int64 `json:"itemid"`
		} `json:"params"`
		Actionlog Actionlog `json:"actionlog"`
		Scheme    string    `json:"scheme"`
	} `json:"buttons,omitempty"`
}

type Actionlog struct {
	ActCode     int         `json:"act_code"`
	Fid         string      `json:"fid"`
	Lfid        string      `json:"lfid"`
	Oid         interface{} `json:"oid"`
	Cardid      string      `json:"cardid"`
	Ext         string      `json:"ext"`
	Featurecode string      `json:"featurecode,omitempty"`
	Mark        string      `json:"mark,omitempty"`
	Uicode      string      `json:"uicode,omitempty"`
	Luicode     string      `json:"luicode,omitempty"`
	Lcardid     string      `json:"lcardid,omitempty"`
}

type Element struct {
	Uid       int       `json:"uid"`
	Scheme    string    `json:"scheme"`
	Itemid    string    `json:"itemid"`
	Actionlog Actionlog `json:"actionlog"`
}
