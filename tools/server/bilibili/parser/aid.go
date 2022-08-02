package parser

import (
	"fmt"
	"github.com/actliboy/hoper/server/go/lib/utils/fs"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client/crawler"
	"github.com/tidwall/gjson"
	"tools/bilibili/api"
	"tools/bilibili/fetcher"
	"tools/bilibili/model"
)

var _getAidUrlTemp = "https://api.bilibili.com/x/space/arc/search?mid=%d&ps=30&tid=0&pn=%d&keyword=&order=pubdate&jsonp=jsonp"

//var _getCidUrlTemp = "https://api.bilibili.com/x/player/pagelist?aid=%d"

func UpSpaceParseFun(contents []byte) ([]*crawler.Request, error) {

	value := gjson.GetManyBytes(contents, "data.list.vlist", "data.page")

	requests, upid := getAidDetailReqList(value[0])
	requests = append(requests, getNewBilibiliUpSpaceReqList(value[1], upid)...)

	return requests, nil

}

func getAidDetailReqList(pageInfo gjson.Result) ([]*crawler.Request, int64) {

	var retRequests []*crawler.Request
	var upid int64
	for _, i := range pageInfo.Array() {
		aid := i.Get("aid").Int()
		upid = i.Get("mid").Int()
		title := i.Get("title").String()
		title = fs.PathClean(title) // remove special characters
		reqUrl := api.GetViewUrl(aid)
		videoAid := model.NewVideoAidInfo(aid, title)
		reqParseFunction := GenGetAidChildrenParseFun(videoAid) //子视频
		req := crawler.NewRequest2(reqUrl, fetcher.DefaultFetcher, reqParseFunction)
		retRequests = append(retRequests, req)
	}
	return retRequests, upid
}

// 访问up主的时候 需要翻页
func getNewBilibiliUpSpaceReqList(pageInfo gjson.Result, upid int64) []*crawler.Request {

	var retRequests []*crawler.Request

	count := pageInfo.Get("count").Int()
	pn := pageInfo.Get("pn").Int()
	ps := pageInfo.Get("ps").Int()
	var extraPage int64
	if count%ps > 0 {
		extraPage = 1
	}
	totalPage := count/ps + extraPage
	for i := int64(1); i <= totalPage; i++ {
		if i == pn {
			continue
		}
		reqUrl := fmt.Sprintf(_getAidUrlTemp, upid, i)
		req := crawler.NewRequest2(reqUrl, fetcher.DefaultFetcher, UpSpaceParseFun)
		retRequests = append(retRequests, req)
	}
	return retRequests
}

func GetRequestByUpId(upid int64) *crawler.Request {

	reqUrl := fmt.Sprintf(_getAidUrlTemp, upid, 1)
	return crawler.NewRequest2(reqUrl, fetcher.DefaultFetcher, UpSpaceParseFun)
}
