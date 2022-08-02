package parser

import (
	"crypto/md5"
	"fmt"
	"github.com/actliboy/hoper/server/go/lib/utils/fs"
	gcrawler "github.com/actliboy/hoper/server/go/lib/utils/generics/net/http/client/crawler"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client/crawler"
	"github.com/tidwall/gjson"
	"tools/bilibili/api"
	"tools/bilibili/download"
	"tools/bilibili/fetcher"
	"tools/bilibili/model"
	"tools/bilibili/tool"

	"strconv"
)

var _entropy = "rbMCKn@KuamXWlPMoJGsKcbiJKUfkPF_8dABscJntvqhRSETg"
var _paramsTemp = "appkey=%s&cid=%s&otype=json&qn=%s&quality=%s&type="
var _playApiTemp = "https://interface.bilibili.com/v2/playurl?%s&sign=%s"
var _quality = "116"

func GenGetAidChildrenParseFun(videoAid *model.VideoAid) crawler.ParseFun {
	return func(contents []byte) ([]*crawler.Request, error) {
		var requests []*crawler.Request
		if videoAid.Title == strconv.FormatInt(videoAid.Aid, 10) { // call from aid-related, we need to get the title of the video
			title := gjson.GetBytes(contents, "data.title").String()
			title = fs.PathClean(title) // remove special characters
			videoAid.Title = title
		}
		data := gjson.GetBytes(contents, "data.pages").Array()
		fmt.Println("即将开始下载：", videoAid.Title)
		appKey, sec := tool.GetAppKey(_entropy)

		var videoTotalPage int64
		for _, i := range data {
			cid := i.Get("cid").Int()
			page := i.Get("page").Int()
			part := i.Get("part").String()
			part = fs.PathClean(part) //remove special characters
			videoCid := model.NewVideoCidInfo(cid, videoAid, page, part)
			videoTotalPage += 1
			cidStr := strconv.FormatInt(videoCid.Cid, 10)

			params := fmt.Sprintf(_paramsTemp, appKey, cidStr, _quality, _quality)
			chksum := fmt.Sprintf("%x", md5.Sum([]byte(params+sec)))

			urlApi := fmt.Sprintf(_playApiTemp, params, chksum)

			req := crawler.NewRequest2(urlApi, fetcher.DefaultFetcher, GenVideoDownloadParseFun(videoCid))
			requests = append(requests, req)
		}

		videoAid.SetPage(videoTotalPage)

		return requests, nil
	}
}

func GetRequestByAid(aid int64) *crawler.Request {
	reqUrl := api.GetViewUrl(aid)
	videoAid := model.NewVideoAidInfo(aid, fmt.Sprintf("%d", aid))
	reqParseFunction := GenGetAidChildrenParseFun(videoAid)
	req := crawler.NewRequest2(reqUrl, fetcher.DefaultFetcher, reqParseFunction)
	return req
}

func GetRequestByFav(aid int64) *crawler.Request {
	return gcrawler.NewRequest(api.GetViewUrl(aid), download.ViewInfoHandleFun)
}

func FavReqs(page int) []*crawler.Request {
	var requests []*crawler.Request
	for i := 1; i <= page; i++ {
		req := gcrawler.NewRequest(api.GetFavListUrl(i), FavList)
		requests = append(requests, req)
	}
	return requests
}

func FavReq(page int) *crawler.Request {
	return gcrawler.NewRequest(api.GetFavListUrl(page), FavList)
}

var apiService = &api.API{}

func FavList(url string) ([]*crawler.Request, error) {
	res, err := api.Get[*api.FavList](url)
	if err != nil {
		return nil, err
	}
	var requests []*crawler.Request
	for _, fav := range res.Medias {
		aid := tool.Bv2av(fav.Bvid)
		req := GetRequestByFav(aid)
		requests = append(requests, req)
	}
	return requests, nil
}
