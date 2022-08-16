package download

import (
	gcrawler "github.com/actliboy/hoper/server/go/lib/utils/generics/net/http/client/crawler"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client/crawler"
	"tools/bilibili/rpc"
	"tools/bilibili/tool"
)

func FavReqs(pageStart, pageEnd int, handleFun crawler.HandleFun) []*crawler.Request {
	var requests []*crawler.Request
	for i := pageStart; i <= pageEnd; i++ {
		req := gcrawler.NewRequest(rpc.GetFavListUrl(i), handleFun)
		requests = append(requests, req)
	}
	return requests
}

func GetByBvId(id string, handleFun crawler.HandleFun) *crawler.Request {
	avid := tool.Bv2av(id)
	return GetViewInfoReq(avid, handleFun)
}

func UpSpaceList(upid int, handleFun crawler.HandleFun) *crawler.Request {
	return gcrawler.NewRequest(rpc.GetUpSpaceListUrl(upid, 1), UpSpaceListFirstPageHandleFun(upid))
}

func GetViewInfoReq(aid int, handleFun crawler.HandleFun) *crawler.Request {
	return crawler.NewKindRequest(rpc.GetViewUrl(aid), KindViewInfo, handleFun)
}
