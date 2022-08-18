package rpc

import (
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

func GetV[T any](url string) T {
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

func (api *API) GetView(aid int) *ViewInfo {
	return GetV[*ViewInfo](GetViewUrl(aid))
}

func (api *API) GetNav() *NavInfo {
	return GetV[*NavInfo](GetNavUrl())
}

func (api *API) GetFavList(page int) *FavList {
	return GetV[*FavList](GetFavListUrl(page))
}

func (api *API) GetPlayerInfo(avid, cid, qn int) *VideoInfo {
	return GetV[*VideoInfo](GetPlayerUrl(avid, cid, qn))
}

func GetViewUrl(aid int) string {
	return fmt.Sprintf("%s/x/web-interface/view?aid=%d", Host, aid)
}

func GetNavUrl() string {
	return fmt.Sprintf("%s/x/web-interface/nav", Host)
}

func GetFavListUrl(page int) string {
	return fmt.Sprintf("%s/x/v3/fav/resource/list?media_id=63181530&pn=%d&ps=20&keyword=&order=mtime&type=0&tid=0&platform=web&jsonp=jsonp", Host, page)
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
