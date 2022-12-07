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

func GetUserInfo(uid int) *User {
	api := CommonCommonIndex + "100505" + strconv.Itoa(uid)
	user, err := Get[User](api)
	if err != nil {
		log.Error(api, err)
	}
	return user
}

func GetUserInfoV2(uid int) *User {
	api := CommonIndex + "type=uid&value=" + strconv.Itoa(uid)
	user, err := Get[User](api)
	if err != nil {
		log.Error(api, err)
	}
	return user
}

func SearchUserWeibo(uid int, q string, page int) *User {
	api := CommonCommonIndex + "100103type=401&q=" + q + "&container_ext=profile_uid:" + strconv.Itoa(uid) + "&page_type=searchall&page=" + strconv.Itoa(page)
	user, err := Get[User](api)
	if err != nil {
		log.Error(err)
	}
	return user
}

func GetUserWeibo(uid int, page int) *WeiboList {
	api := CommonCommonIndex + "230413" + strconv.Itoa(uid) + "&page=" + strconv.Itoa(page)
	weibo, err := Get[WeiboList](api)
	if err != nil {
		log.Error(api, err)
	}
	return weibo
}

func GetLongWeibo(uid int) *User {
	api := Host + "/detail/" + strconv.Itoa(uid)
	user, err := Get[User](api)
	if err != nil {
		log.Error(api, err)
	}
	return user
}

func GetLivePhoto() {
	const prefix = "https://video.weibo.com/media/play?livephoto=//us.sinaimg.cn/"
}

func GetComments() {
	const api = Host + "/api/comments/show?id="
}

func GetFollows(uid, page int) (*Follow, error) {
	api := fmt.Sprintf(CommonCommonIndex+"231051_-_followers_-_%d&page=%d", uid, page)
	follow, err := Get[Follow](api)
	if err != nil {
		log.Error(api, err)
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
		log.Error(api, err)
		return nil, err
	}
	return cards, nil
}
