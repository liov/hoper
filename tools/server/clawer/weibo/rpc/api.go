package rpc

import (
	"fmt"
	"github.com/liov/hoper/server/go/lib/utils/log"
	httpi "github.com/liov/hoper/server/go/lib/utils/net/http"
	"github.com/liov/hoper/server/go/lib/utils/net/http/client"
	"strconv"
)

const Host = "https://m.weibo.cn"
const CommonIndex = Host + "/api/container/getIndex?"
const CommonSecond = Host + "/api/container/getSecond?"
const CommonCommonIndex = CommonIndex + "containerid="
const CommonCommonSecond = CommonSecond + "containerid="

var Cookie = ``

var commonReq, commonReqWithoutCookie *client.RequestParams

func GetReqCommon() *client.RequestParams {
	if commonReq != nil {
		return commonReq
	}
	commonReq = client.NewGetRequest("").AddHeader(httpi.HeaderCookie, Cookie).DisableLog().CachedHeader("GetReqCommon")
	return commonReq
}

func GetCommonReqWithoutCookie() *client.RequestParams {
	if commonReqWithoutCookie != nil {
		return commonReqWithoutCookie
	}
	commonReqWithoutCookie = client.NewGetRequest("").DisableLog().CachedHeader("GetCommonReqWithoutCookie")
	return commonReqWithoutCookie
}

type Response[T any] struct {
	Ok   int `json:"ok"`
	Data *T  `json:"data"`
}

func Get[T any](url string) (*T, error) {
	var res Response[T]
	err := GetReqCommon().CacheGet(url, &res)
	if err != nil {
		return new(T), err
	}
	return res.Data, nil
}

func GetWithoutCookie[T any](url string) (*T, error) {
	var res Response[T]
	err := GetCommonReqWithoutCookie().CacheGet(url, &res)
	if err != nil {
		return new(T), err
	}
	return res.Data, nil
}

func GetUserInfo(uid int) (*User, error) {
	api := CommonCommonIndex + "100505" + strconv.Itoa(uid)
	user, err := Get[User](api)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserInfoV2(uid int) (*User, error) {
	api := CommonIndex + "type=uid&value=" + strconv.Itoa(uid)
	user, err := Get[User](api)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func SearchUserWeibo(uid int, q string, page int) (*WeiboList, error) {
	api := CommonCommonIndex + "100103type=401&q=" + q + "&container_ext=profile_uid:" + strconv.Itoa(uid) + "&page_type=searchall&page=" + strconv.Itoa(page)
	user, err := Get[WeiboList](api)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserWeibo(uid int, page int) *WeiboList {
	api := CommonCommonIndex + "230413" + strconv.Itoa(uid) + "&page=" + strconv.Itoa(page)
	weibo, err := Get[WeiboList](api)
	if err != nil {
		log.Error(api, err)
	}
	return weibo
}

func GetLongWeibo(wid string) (*Mblog, error) {
	api := Host + "/statuses/show?id=" + wid
	mblog, err := Get[Mblog](api)
	if err != nil {
		return nil, err
	}
	return mblog, nil
}

func GetLivePhoto() {
	const prefix = "https://video.weibo.com/media/play?livephoto=//us.sinaimg.cn/"
}

func GetComments(wid string, page int) (*CommentList, error) {
	api := Host + "/api/comments/show?id=" + wid + "&page=" + strconv.Itoa(page)
	comment, err := Get[CommentList](api)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func GetFollows(uid, page int) (*Follow, error) {
	api := fmt.Sprintf(CommonCommonIndex+"231051_-_followers_-_%d&page=%d", uid, page)
	follow, err := Get[Follow](api)
	if err != nil {
		return nil, err
	}
	return follow, nil
}

func GetPicCards(uid int) *PicCards {
	api := fmt.Sprintf(CommonCommonIndex+"107803%d", uid)
	cards, err := Get[PicCards](api)
	if err != nil {
		log.Error(api, err)
	}
	return cards
}

func DownloadPic(url string) {

}

func GetPhotos(uid, page int) (*PicCards, error) {
	api := fmt.Sprintf(CommonCommonSecond+"107803%d_-_photoall&page=%d&count=20", uid, page)
	cards, err := Get[PicCards](api)
	if err != nil {
		return nil, err
	}
	return cards, nil
}

func GetVideos(uid, page int) (*Weibo, error) {
	api := fmt.Sprintf(CommonCommonIndex+"230413%d_-_WEIBO_SECOND_PROFILE_WEIBO_VIDEO&page=%d&count=20", uid, page)
	cards, err := Get[Weibo](api)
	if err != nil {
		return nil, err
	}
	return cards, nil
}

func GetChannels(channel ChannelType, uid, page int) (*Weibo, error) {
	api := fmt.Sprintf(CommonCommonIndex+"230413%d_-_WEIBO_SECOND_PROFILE_WEIBO%s&page=%d&count=20", uid, channel, page)
	cards, err := Get[Weibo](api)
	if err != nil {
		return nil, err
	}
	return cards, nil
}

// 部分图片返回格式不规范
func GetChannelsV2(channel ChannelType, uid, page int) (*WeiboV2, error) {
	api := fmt.Sprintf(CommonCommonIndex+"230413%d_-_WEIBO_SECOND_PROFILE_WEIBO%s&page=%d&count=20", uid, channel, page)
	cards, err := Get[WeiboV2](api)
	if err != nil {
		return nil, err
	}
	return cards, nil
}

// 全部
const ChannelApi = CommonCommonIndex + "230413%d_-_WEIBO_SECOND_PROFILE_WEIBO%s&page=%d&count=20"

type ChannelType string

const (
	ALL     ChannelType = ""         // All
	ORI     ChannelType = "_ORI"     // 原创
	VIDEO   ChannelType = "_VIDEO"   // 视频
	ARTICAL ChannelType = "_ARTICAL" // 文章
	PIC     ChannelType = "_PIC"     // 图片
)

// 部分图片返回格式不规范
func GetFollowsWeibo(maxId string) (*WeiboFollowsList, error) {
	api := Host + "/feed/friends"
	if maxId != "" {
		api += "?max_id=" + maxId
	}
	list, err := Get[WeiboFollowsList](api)
	if err != nil {
		return nil, err
	}
	return list, nil
}
