package download

import (
	"github.com/actliboy/hoper/server/go/lib/utils/fs"
	gcrawler "github.com/actliboy/hoper/server/go/lib/utils/generics/net/http/client/crawler"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client/crawler"
	"tools/bilibili/dao"
	"tools/bilibili/rpc"
	"tools/bilibili/tool"
)

func FavReqs(pageBegin, pageEnd int, handleFun crawler.HandleFun) []*crawler.Request {
	if pageEnd < pageBegin {
		pageEnd = pageBegin
	}
	var requests []*crawler.Request
	for i := pageBegin; i <= pageEnd; i++ {
		req := gcrawler.NewRequest(rpc.GetFavListUrl(63181530, i), handleFun)
		requests = append(requests, req)
	}
	return requests
}

func FavVideo(engine *crawler.Engine) {

	for {
		var videos []*Video
		dao.Dao.Hoper.DB.Raw(`SELECT
    a.aid,b.cid,a.data->'title' title,
    p->'page' page,p->'part' part
FROM
    "bilibili"."view" a,jsonb_path_query(a.data,'$.pages[*]') AS p
LEFT JOIN "bilibili"."video" b ON (p->'cid')::int8 = b.cid
WHERE b.record = false
LIMIT 20;`).Find(&videos)
		for _, video := range videos {
			video.Title = fs.PathClean(video.Title)
			req := crawler.NewKindRequest(rpc.GetPlayerUrl(video.Aid, video.Cid, 120), KindGetPlayerUrl, video.PlayerUrlHandleFun)
			engine.Engine.AddTask(engine.NewTask(req))
		}
	}
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
