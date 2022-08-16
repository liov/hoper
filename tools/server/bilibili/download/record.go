package download

import (
	"context"
	"encoding/json"
	"github.com/actliboy/hoper/server/go/lib/utils/dao/db/postgres"
	"github.com/actliboy/hoper/server/go/lib/utils/fs"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client/crawler"
	"tools/bilibili/dao"
	"tools/bilibili/rpc"
	"tools/bilibili/tool"
)

func RecordFavList(ctx context.Context, url string) ([]*crawler.Request, error) {
	res, err := rpc.Get[*rpc.FavList](url)
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
			req1 := GetViewInfoReq(aid, RecordViewInfoHandleFun)
			req2 := crawler.NewKindRequest(fav.Cover, KindDownloadCover, DownloadCover(ctx, fav.Id))
			requests = append(requests, req1, req2)
		}
	}
	return requests, nil
}

func RecordViewInfoHandleFun(ctx context.Context, url string) ([]*crawler.Request, error) {
	res, err := rpc.Get[rpc.ViewInfo](url)
	if err != nil {
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
		video := &Video{fs.PathClean(res.Title), res.Aid, page.Cid, page.Page, page.Part, 0}

		exists, err := bilibiliDao.VideoExists(video.Aid, video.Cid)
		if err != nil {
			return nil, err
		}
		if !exists {
			req := crawler.NewKindRequest(rpc.GetPlayerUrl(res.Aid, page.Cid, 120), KindGetPlayerUrl, video.RecordDownloadHandleFun)
			requests = append(requests, req)
		}
	}
	return requests, nil
}

func (video *Video) RecordDownloadHandleFun(ctx context.Context, url string) ([]*crawler.Request, error) {
	res, err := rpc.Get[*rpc.VideoInfo](url)
	if err != nil {
		return nil, err
	}

	video.Quality = res.Quality
	var requests []*crawler.Request
	for _, durl := range res.Durl {
		req := crawler.NewKindRequest(durl.Url, KindDownloadVideo, video.GetDownloadHandleFun(durl.Order))
		requests = append(requests, req)
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

	return requests, nil
}
