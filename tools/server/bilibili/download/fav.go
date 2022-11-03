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

	for _, favId := range favIds {
		go func(favId int) {
			page := 1
			favIdStr := strconv.Itoa(favId)
			cancel := make(chan struct{}, 1)
		Loop:
			for {
				select {
				case <-timer.C:
					taskFun := RecordFavListReqAfterRecordView(favId, page, lastRecordTime, cancel)
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

func RecordFavListReqAfterRecordView(favId, page int, lastRecordTime time.Time, cancel chan struct{}) crawler.TaskFunc {
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
				if !strings.HasSuffix(fav.Cover, "be27fd62c99036dce67efface486fb0a88ffed06.jpg") {
					req1 := RecordViewInfoReqAfterRecordVideo(aid)
					req2 := CoverDownloadReq(fav.Cover, fav.Upper.Mid, fav.Id)
					requests = append(requests, req1, req2)
				}
			} else if createdAt.Before(lastRecordTime) {
				cancel <- struct{}{}
				return requests, nil
			}
		}
		return requests, nil
	}
}

func RecordFavList(ctx context.Context, url string) ([]*crawler.Request, error) {
	res, err := rpc.Get[*rpc.FavResourceList](url)
	if err != nil {
		return nil, err
	}
	var requests []*crawler.Request
	for _, fav := range res.Medias {
		aid := tool.Bv2av(fav.Bvid)
		bilibiliDao := dao.NewDao(ctx, dao.Dao.Hoper.DB)
		exists, err := bilibiliDao.ViewExists(aid)
		if err != nil {
			return nil, err
		}
		if !exists {
			if !strings.HasSuffix(fav.Cover, "be27fd62c99036dce67efface486fb0a88ffed06.jpg") {
				req1 := RecordViewInfoReq(aid)
				req2 := CoverDownloadReq(fav.Cover, fav.Upper.Mid, fav.Id)
				requests = append(requests, req1, req2)
			}
		}
	}
	return requests, nil
}
