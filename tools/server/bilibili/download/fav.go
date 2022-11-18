package download

import (
	"context"
	"github.com/actliboy/hoper/server/go/lib/utils/generics/net/http/client/crawler"

	"log"
	"strconv"
	"strings"
	"time"
	"tools/bilibili/dao"
	"tools/bilibili/rpc"
	"tools/bilibili/tool"
)

var apiservice = &rpc.API{}

func RecordFavTimer(ctx context.Context, engine *crawler.Engine) {
	favIds := []int{63181530, 62504730}
	timer := time.NewTicker(time.Second)
	lastRecordTime, err := dao.NewDao(ctx, dao.Dao.Hoper.DB).LastCreated(dao.TableNameView)
	if err != nil {
		log.Println(err)
		return
	}
	cancel := make(chan struct{}, 1)

	for _, favId := range favIds {
		go func(favId int) {
			page := 1
			favIdStr := strconv.Itoa(favId)
		Loop:
			for {
				select {
				case <-timer.C:
					taskFun := GetFavListReqAfterRecordView(favId, page, lastRecordTime, cancel)
					engine.AddTask(engine.NewTask(crawler.NewRequest(favIdStr+strconv.Itoa(page), KindRecordFavList, taskFun)))
				case <-cancel:
					timer.Stop()
					break Loop
				}
				page++
			}
		}(favId)
	}
}

func GetFavListReqAfterRecordView(favId, page int, lastRecordTime time.Time, cancel chan struct{}) crawler.TaskFunc {
	return func(ctx context.Context) ([]*crawler.Request, error) {
		res, err := apiservice.GetFavLResourceList(favId, page)
		if err != nil {
			return nil, err
		}
		zeroTime := time.Time{}
		var requests []*crawler.Request
		for _, fav := range res.Medias {
			aid := tool.Bv2av(fav.Bvid)
			bilibiliDao := dao.NewDao(ctx, dao.Dao.Hoper.DB)
			createdAt, err := bilibiliDao.ViewCreatedTime(aid)
			if err != nil {
				return nil, err
			}
			if createdAt == zeroTime {
				if !strings.HasSuffix(fav.Cover, NULLCOVER) {
					req1 := RecordViewInfoReqAfterRecordVideo(aid)
					req2 := CoverDownloadReq(fav.Cover, fav.Upper.Mid, fav.Id)
					requests = append(requests, req1, req2)
				}
			} else if createdAt.Before(lastRecordTime) {
				close(cancel)
				return requests, nil
			}
		}
		return requests, nil
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
					//req2 := CoverDownloadReq(fav.Cover, fav.Upper.Mid, fav.Id)
					requests = append(requests, req1)
				}

			}
			return requests, nil
		},
	}
}

func FixRecordFav(ctx context.Context, engine *crawler.Engine) {
	for page := 2; page < 20; page++ {
		log.Printf("第%d页", page)
		engine.AddTask(engine.NewTask(RecordFavReq(63181530, page)))
		page++
	}
	for page := 1; page < 10; page++ {
		log.Printf("第%d页", page)
		engine.AddTask(engine.NewTask(RecordFavReq(62504730, page)))
		page++
	}
}
