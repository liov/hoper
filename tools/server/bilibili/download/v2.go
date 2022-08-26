package download

import (
	"context"
	"encoding/json"
	"github.com/actliboy/hoper/server/go/lib/utils/conctrl"
	"github.com/actliboy/hoper/server/go/lib/utils/dao/db/postgres"
	"github.com/actliboy/hoper/server/go/lib/utils/fs"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client/crawler"
	"gorm.io/gorm"
	"tools/bilibili/dao"
	"tools/bilibili/rpc"
	"tools/bilibili/tool"
)

var apiservice = &rpc.API{}

func FavListV2(ctx context.Context) ([]*crawler.Request, error) {
	res, err := apiservice.GetFavLResourceList(1, 20)
	if err != nil {
		return nil, err
	}
	var requests []*crawler.Request
	for _, fav := range res.Medias {
		aid := tool.Bv2av(fav.Bvid)
		req1 := GetViewInfoReqV2(aid)
		req2 := crawler.NewUrlKindRequest(fav.Cover, KindDownloadCover, CoverDownload(ctx, fav.Id))
		requests = append(requests, req1, req2)
	}
	return requests, nil
}

func GetViewInfoReqV2(aid int) *crawler.Request {
	return &crawler.Request{
		TaskMeta: conctrl.TaskMeta{Kind: KindViewInfo},
		TaskFun: func(ctx context.Context) ([]*crawler.Request, error) {
			view, err := apiservice.GetView(aid)
			if err != nil || view.Aid == 0 {
				return nil, err
			}
			bilibiliDao := dao.NewDao(ctx, dao.Dao.Hoper.DB)
			exists, err := bilibiliDao.ViewExists(view.Aid)
			if err != nil {
				return nil, err
			}
			if !exists {
				data, err := json.Marshal(view)
				if err != nil {
					return nil, err
				}
				err = bilibiliDao.CreateView(&dao.View{
					Bvid:        view.Bvid,
					Aid:         view.Aid,
					Data:        data,
					CoverRecord: false,
				})
				if err != nil && !postgres.IsDuplicate(err) {
					return nil, err
				}
			}
			var requests []*crawler.Request
			for _, page := range view.Pages {
				video := &Video{fs.PathClean(view.Title), view.Aid, page.Cid, page.Page, page.Part, 0}

				req := crawler.NewKindRequest("", KindGetPlayerUrl, video.PlayerUrlHandleFunV2)
				requests = append(requests, req)
			}
			return requests, nil
		},
	}
}

func (video *Video) PlayerUrlHandleFunV2(ctx context.Context) ([]*crawler.Request, error) {
	var dvideo dao.Video
	err := dao.Dao.Hoper.Table(dao.TableNameVideo).Select("cid,record").Where("cid = ?", video.Cid).First(&dvideo).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if dvideo.Record {
		return nil, nil
	}
	res, err := apiservice.GetPlayerInfo(video.Aid, video.Cid, 120)
	if err != nil {
		return nil, err
	}

	video.Quality = res.Quality
	var requests []*crawler.Request
	if !dvideo.Record {
		for _, durl := range res.Durl {
			err = video.DownloadVideoHandleFun(durl.Order, durl.Url)
			if err != nil {
				return nil, err
			}
		}
	}

	bilibiliDao := dao.NewDao(ctx, dao.Dao.Hoper.DB)
	if err == gorm.ErrRecordNotFound || dvideo.Cid == 0 {
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
	}

	return requests, nil
}
