package download

import (
	"github.com/actliboy/hoper/server/go/lib/utils/dao/db/postgres"
	"github.com/actliboy/hoper/server/go/lib/utils/fs"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client/crawler"
	"math"
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
		req1 := crawler.NewUrlRequest(rpc.GetFavResourceListUrl(63181530, i), handleFun) //62504730,63181530
		req2 := crawler.NewUrlRequest(rpc.GetFavResourceListUrl(62504730, i), handleFun) //62504730,63181530
		requests = append(requests, req1, req2)
	}
	return requests
}

func FavVideo(engine *crawler.Engine) {
	minAid := math.MaxInt
	for {
		var videos []*Video
		dao.Dao.Hoper.DB.Raw(`SELECT b.aid,b.cid,a.title,a.p->'page' page,a.p->'part' part
FROM `+dao.TableNameVideo+` b 
LEFT JOIN (SELECT data->'title' title ,jsonb_path_query(data,'$.pages[*]') p FROM `+dao.TableNameView+`)  a ON (a.p->'cid')::int8 = b.cid
WHERE b.record = false AND b.aid < ? AND b.`+postgres.NotDeleted+` ORDER BY b.aid DESC LIMIT 20`, minAid).Find(&videos)
		if len(videos) == 0 {
			return
		}
		for _, video := range videos {
			if video.Title == "" {
				req := ViewRecordUpdate(video.Aid)
				engine.Engine.AddTask(engine.NewTask(req))
			} else {
				video.Title = fs.PathClean(video.Title)
				req := crawler.NewUrlKindRequest(rpc.GetPlayerUrl(video.Aid, video.Cid, 120), KindGetPlayerUrl, video.PlayerUrlHandleFun)
				engine.Engine.AddTask(engine.NewTask(req))
			}
		}
		minAid = videos[len(videos)-1].Aid
	}
}

func ViewUpdate(aid int) *crawler.Request {
	video := &Video{Aid: aid, Cid: 0}
	return crawler.NewUrlKindRequest(rpc.GetPlayerUrl(aid, 0, 120), KindGetPlayerUrl, video.VideoRecord)
}

func GetByBvId(id string, handleFun crawler.HandleFun) *crawler.Request {
	avid := tool.Bv2av(id)
	return GetViewInfoReq(avid, handleFun)
}

func UpSpaceList(upid int, handleFun crawler.HandleFun) *crawler.Request {
	return crawler.NewUrlRequest(rpc.GetUpSpaceListUrl(upid, 1), UpSpaceListFirstPageHandleFun(upid))
}

func GetViewInfoReq(aid int, handleFun crawler.HandleFun) *crawler.Request {
	return crawler.NewUrlKindRequest(rpc.GetViewUrl(aid), KindViewInfo, handleFun)
}
