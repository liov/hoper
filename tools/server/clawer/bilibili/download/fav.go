package download

import (
	"context"
	"github.com/liov/hoper/server/go/lib_v2/utils/net/http/client/crawler"
	"log"
	"strconv"
	"strings"
	"time"
	"tools/clawer/bilibili/dao"
	"tools/clawer/bilibili/rpc"
	"tools/clawer/bilibili/tool"
)

var apiservice = &rpc.API{}

func RecordFavTimer(startTime time.Time) []*crawler.Request {
	favIds := []int{62504730, 63181530}
	var requests []*crawler.Request
	for _, favId := range favIds {
		requests = append(requests, GetFavListReqAfterRecordView(favId, 1, startTime))
	}
	return requests
}

var zeroTime = time.Time{}

func GetFavListReqAfterRecordView(favId, page int, startTime time.Time) *crawler.Request {
	return &crawler.Request{
		TaskMeta: crawler.TaskMeta{BaseTaskMeta: crawler.BaseTaskMeta{Key: strconv.Itoa(favId) + strconv.Itoa(page)}, Kind: KindRecordFavList},
		TaskFunc: func(ctx context.Context) ([]*crawler.Request, error) {
			log.Printf("获取收藏夹%d,第%d页\n", favId, page)
			res, err := apiservice.GetFavLResourceList(favId, page)
			if err != nil {
				return nil, err
			}

			var requests []*crawler.Request
			for _, fav := range res.Medias {
				aid := tool.Bv2av(fav.Bvid)
				bilibiliDao := dao.NewDao(ctx, dao.Dao.Hoper.DB)
				createdTime, err := bilibiliDao.ViewCreatedTime(aid)
				if err != nil {
					return nil, err
				}
				if createdTime == zeroTime {
					if !strings.HasSuffix(fav.Cover, NULLCOVER) {
						req1 := RecordViewInfoReqAfterRecordVideo(aid)
						req2 := CoverDownloadReq(fav.Pubtime, fav.Cover, fav.Upper.Mid, fav.Id)
						requests = append(requests, req1, req2)
					}
				} else if createdTime.Before(startTime) {
					return requests, nil
				}

			}
			requests = append(requests, GetFavListReqAfterRecordView(favId, page+1, startTime))
			return requests, nil
		},
	}
}

func RecordFavReq(favId, page int) *crawler.Request {
	return &crawler.Request{
		TaskMeta: crawler.TaskMeta{BaseTaskMeta: crawler.BaseTaskMeta{Key: strconv.Itoa(favId) + strconv.Itoa(page)}, Kind: KindRecordFavList},
		TaskFunc: func(ctx context.Context) ([]*crawler.Request, error) {
			res, err := apiservice.GetFavLResourceList(favId, page)
			if err != nil {
				return nil, err
			}
			var requests []*crawler.Request
			for _, fav := range res.Medias {
				aid := tool.Bv2av(fav.Bvid)
				if !strings.HasSuffix(fav.Cover, NULLCOVER) {
					req1 := RecordViewInfoReqAfterRecordVideo(aid)
					req2 := CoverDownloadReq(fav.Pubtime, fav.Cover, fav.Upper.Mid, fav.Id)
					requests = append(requests, req1, req2)
				}

			}
			return requests, nil
		},
	}
}

func FixRecordFav(engine *crawler.Engine) {
	for page := 1; page < 5; page++ {
		log.Printf("第%d页", page)
		engine.AddTask(engine.BaseTask(RecordFavReq(63181530, page)))
		page++
	}
	for page := 1; page < 2; page++ {
		log.Printf("第%d页", page)
		engine.AddTask(engine.BaseTask(RecordFavReq(62504730, page)))
		page++
	}
}
