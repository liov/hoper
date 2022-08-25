package rpc

import (
	"errors"
	"fmt"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	httpi "github.com/actliboy/hoper/server/go/lib/utils/net/http"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client"
)

type API struct{}

var api = &API{}

const Host = "https://api.bilibili.com"

var Cookie = ``

func ReqCommonSet(req *client.RequestParams) *client.RequestParams {
	return ReqCommonSetWithoutCookie(req).AddHeader(httpi.HeaderCookie, Cookie)
}

func ReqCommonSetWithoutCookie(req *client.RequestParams) *client.RequestParams {
	return req.AddHeader(httpi.HeaderUserAgent, client.UserAgent1).AddHeader(httpi.HeaderCookie, Cookie).DisableLog()
}

func Get[T any](url string) (T, error) {
	var res Response[T]
	err := ReqCommonSet(client.NewGetRequest(url)).Do(nil, &res)
	if err != nil {
		return *new(T), err
	}
	return res.Data, nil
}

func GetWithoutCookie[T any](url string) (T, error) {
	var res Response[T]
	err := ReqCommonSetWithoutCookie(client.NewGetRequest(url)).Do(nil, &res)
	if err != nil {
		return *new(T), err
	}
	return res.Data, nil
}

func GetIgnoreErr[T any](url string) T {
	var res Response[T]
	err := ReqCommonSet(client.NewGetRequest(url)).Do(nil, &res)
	if err != nil {
		log.Error(err)
	}
	return res.Data
}

func GetZ[T any](url string) *T {
	res := new(T)
	err := ReqCommonSet(client.NewGetRequest(url)).Do(nil, res)
	if err != nil {
		log.Error(err)
	}
	return res
}

func (api *API) GetView(aid int) (*ViewInfo, error) {
	return Get[*ViewInfo](GetViewUrl(aid))
}

func (api *API) GetNav() *NavInfo {
	return GetIgnoreErr[*NavInfo](GetNavUrl())
}

func (api *API) GetFavLResourceList(favId, page int) (*FavResourceList, error) {
	return Get[*FavResourceList](GetFavResourceListUrl(favId, page))
}

func (api *API) GetPlayerInfo(avid, cid, qn int) (*VideoInfo, error) {
	return Get[*VideoInfo](GetPlayerUrl(avid, cid, qn))
}

func GetViewUrl(aid int) string {
	return fmt.Sprintf("%s/x/web-interface/view?aid=%d", Host, aid)
}

func GetNavUrl() string {
	return fmt.Sprintf("%s/x/web-interface/nav", Host)
}

func GetFavResourceListUrl(favId, page int) string {
	return fmt.Sprintf("%s/x/v3/fav/resource/list?media_id=%d&pn=%d&ps=20&keyword=&order=mtime&type=0&tid=0&platform=web&jsonp=jsonp", Host, favId, page)
}

func GetPlayerUrl(avid, cid, qn int) string {
	return fmt.Sprintf("%s/x/player/playurl?avid=%d&cid=%d&qn=%d&fourk=1", Host, avid, cid, qn)
}

func GetUpSpaceListUrl(upid, page int) string {
	var _getAidUrlTemp = "%s/x/space/arc/search?mid=%d&ps=30&tid=0&pn=%d&keyword=&order=pubdate&jsonp=jsonp"
	return fmt.Sprintf(_getAidUrlTemp, Host, upid, page)
}

func (api *API) GetUpSpaceList(upid, page int) (*UpSpaceList, error) {
	return Get[*UpSpaceList](GetUpSpaceListUrl(upid, page))
}

type Response[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ttl     int    `json:"ttl"`
	Data    T      `json:"data"`
}

func (receiver *Response[T]) CheckError() error {
	if receiver.Code != 0 {
		return errors.New(receiver.Message)
	}
	return nil
}

// 收藏夹
func GetFavListUrl(mid int) string {
	var _getAidUrlTemp = "%s/x/v3/fav/folder/created/list-all?up_mid=%d&jsonp=jsonp"
	return fmt.Sprintf(_getAidUrlTemp, Host, mid)
}

func GetFavList(mid int) *FavList {
	return GetIgnoreErr[*FavList](GetFavListUrl(mid))
}

// 收藏订阅
func GetFavCollectedListUrl(mid, page int) string {
	var _getAidUrlTemp = "%s/x/v3/fav/folder/collected/list?pn=%d&ps=20&up_mid=%d&platform=web&jsonp=jsonp"
	return fmt.Sprintf(_getAidUrlTemp, Host, page, mid)
}

func GetFavCollectedList(mid, page int) string {
	var _getAidUrlTemp = "%s/x/v3/fav/folder/collected/list?pn=%d&ps=20&up_mid=%d&platform=web&jsonp=jsonp"
	return fmt.Sprintf(_getAidUrlTemp, Host, page, mid)
}

func GetFavSeasonListUrl(seasonId, page int) string {
	var _getAidUrlTemp = "%sx/space/fav/season/list?season_id=%d&pn=%d&ps=20&jsonp=jsonp"
	return fmt.Sprintf(_getAidUrlTemp, Host, seasonId, page)
}
