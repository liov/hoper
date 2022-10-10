package download

import (
	"context"
	"encoding/json"
	"github.com/actliboy/hoper/server/go/lib/utils/dao/db/postgres"
	"github.com/actliboy/hoper/server/go/lib/utils/fs"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client/crawler"
	"log"
	"strconv"
	"strings"
	"time"
	"tools/bilibili/dao"
	"tools/bilibili/rpc"
	"tools/bilibili/tool"
)

func RecordFav(ctx context.Context, engine *crawler.Engine) {
	favIds := []int{62504730, 63181530}
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
					taskFun := RecordFavList2(favId, page, lastRecordTime, cancel)
					engine.AddTask(engine.NewTask(crawler.NewRequest(favIdStr+strconv.Itoa(page), taskFun)))
				case <-cancel:
					timer.Stop()
					break Loop
				}
				page++
			}
		}(favId)
	}
}

func RecordFavList2(favId, page int, lastRecordTime time.Time, cancel chan struct{}) crawler.TaskFun {
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
					req1 := GetViewInfoReq(aid, ViewInfoRecord)
					req2 := crawler.NewUrlKindRequest(fav.Cover, KindDownloadCover, CoverDownload(ctx, fav.Upper.Mid, fav.Id))
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
				req1 := GetViewInfoReq(aid, ViewInfoRecord)
				requests = append(requests, req1)
				req2 := crawler.NewUrlKindRequest(fav.Cover, KindDownloadCover, CoverDownload(ctx, fav.Upper.Mid, fav.Id))
				requests = append(requests, req2)
			}
		}
	}
	return requests, nil
}

func ViewInfoRecord(ctx context.Context, url string) ([]*crawler.Request, error) {
	res, err := rpc.Get[rpc.ViewInfo](url)
	if err != nil && err.Error() != rpc.ErrorNotFound && err.Error() != rpc.ErrorNotPermission {
		return nil, err
	}
	bilibiliDao := dao.NewDao(ctx, dao.Dao.Hoper.DB)

	data, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	err = bilibiliDao.CreateView(&dao.View{
		Bvid:        res.Bvid,
		Aid:         res.Aid,
		Data:        data,
		CoverRecord: false,
	})
	if err != nil && !postgres.IsDuplicate(err) {
		return nil, err
	}

	var requests []*crawler.Request
	for _, page := range res.Pages {
		video := &Video{res.Owner.Mid, fs.PathClean(res.Title), res.Aid, page.Cid, page.Page, page.Part, 0}

		exists, err := bilibiliDao.VideoExists(video.Cid)
		if err != nil {
			return nil, err
		}
		if !exists {
			req := crawler.NewUrlKindRequest(rpc.GetPlayerUrl(res.Aid, page.Cid, 120), KindGetPlayerUrl, video.VideoRecord)
			requests = append(requests, req)
		}
	}
	return requests, nil
}

func (video *Video) VideoRecord(ctx context.Context, url string) ([]*crawler.Request, error) {
	res, err := rpc.GetWithoutCookie[*rpc.VideoInfo](url)
	if err != nil {
		if err.Error() == rpc.ErrorNotFound {
			dao.Dao.Hoper.Table(dao.TableNameVideo).Where(`cid = ?`, video.Cid).UpdateColumn("deleted_at", time.Now())
			return nil, nil
		}
		return nil, err
	}

	bilibiliDao := dao.NewDao(ctx, dao.Dao.Hoper.DB)
	res.JsonClean()
	data, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	err = bilibiliDao.CreateVideo(&dao.Video{
		Aid:    video.Aid,
		Cid:    video.Cid,
		Data:   data,
		Record: false,
	})
	if err != nil && !postgres.IsDuplicate(err) {
		return nil, err
	}

	return nil, nil
}

func ViewRecordUpdate(aid int) *crawler.Request {
	return crawler.NewRequest("", func(ctx context.Context) ([]*crawler.Request, error) {
		bilibiliDao := dao.NewDao(ctx, dao.Dao.Hoper.DB)
		exists, err := bilibiliDao.ViewExists(aid)
		if err != nil {
			return nil, err
		}
		if !exists {
			req1 := GetViewInfoReq(aid, ViewInfoRecord)
			return []*crawler.Request{req1}, nil
		}

		err = dao.Dao.Hoper.Exec(`INSERT INTO `+dao.TableNameViewBak+`(aid,data) (SELECT aid,data FROM `+dao.TableNameView+` WHERE aid = ?) `, aid).Error
		/*	if err != nil {
			return nil, err
		}*/
		res, err := apiservice.GetView(aid)
		if err != nil || res.Aid == 0 {
			return nil, err
		}
		data, err := json.Marshal(res)
		if err != nil {
			return nil, err
		}
		err = dao.Dao.Hoper.Table(dao.TableNameView).Where(`aid = ?`, aid).Update("data", data).Error
		if err != nil {
			return nil, err
		}
		var requests []*crawler.Request
		for _, page := range res.Pages {
			video := &Video{res.Owner.Mid, fs.PathClean(res.Title), res.Aid, page.Cid, page.Page, page.Part, 0}

			req := crawler.NewUrlKindRequest(rpc.GetPlayerUrl(res.Aid, page.Cid, 120), KindGetPlayerUrl, video.VideoRecord)
			requests = append(requests, req)
		}
		return requests, nil
	})
}
