package download

import (
	"context"
	"encoding/json"

	"github.com/actliboy/hoper/server/go/lib/utils/generics/net/http/client/crawler"

	"github.com/actliboy/hoper/server/go/lib/utils/dao/db/postgres"

	"time"
	"tools/bilibili/dao"
	"tools/bilibili/rpc"
	"tools/bilibili/tool"
)

func RecordViewInfoReqAfterRecordVideo(aid int) *crawler.Request {
	return &crawler.Request{
		TaskMeta: crawler.TaskMeta{Kind: KindViewInfo},
		TaskFunc: func(ctx context.Context) ([]*crawler.Request, error) {
			view, err := RecordViewInfo(ctx, aid)
			if err != nil {
				return nil, err
			}
			var requests []*crawler.Request
			for _, page := range view.Pages {
				if len(view.Pages) == 1 {
					page.Part = PartEqTitle
				}
				video := NewVideo(view.Owner.Mid, view.Title, view.Aid, page.Cid, page.Page, page.Part, time.Now())

				req := video.RecordVideoReqAfterDownloadVideo()
				requests = append(requests, req)
			}
			return requests, nil
		},
	}
}

func RecordViewInfoReq(aid int) *crawler.Request {
	return &crawler.Request{
		TaskMeta: crawler.TaskMeta{Kind: KindViewInfo},
		TaskFunc: func(ctx context.Context) ([]*crawler.Request, error) {
			_, err := RecordViewInfo(ctx, aid)
			return nil, err
		},
	}
}

func RecordViewInfo(ctx context.Context, aid int) (*rpc.ViewInfo, error) {
	view, err := apiservice.GetView(aid)
	if err != nil && err.Error() != rpc.ErrorNotFound && err.Error() != rpc.ErrorNotPermission {
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
	return view, err
}

func ViewRecordUpdateReqAfterRecordVideo(aid int) *crawler.Request {
	return &crawler.Request{
		TaskMeta: crawler.TaskMeta{Kind: KindViewInfo},
		TaskFunc: func(ctx context.Context) ([]*crawler.Request, error) {
			bilibiliDao := dao.NewDao(ctx, dao.Dao.Hoper.DB)
			exists, err := bilibiliDao.ViewExists(aid)
			if err != nil {
				return nil, err
			}
			if !exists {
				req1 := RecordViewInfoReqAfterRecordVideo(aid)
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
				video := NewVideo(res.Owner.Mid, res.Title, res.Aid, page.Cid, page.Page, page.Part, time.Now())

				req := video.RecordVideoReq()
				requests = append(requests, req)
			}
			return requests, nil
		},
	}
}

func RecordViewInfoByBvId(ctx context.Context, id string) (*rpc.ViewInfo, error) {
	avid := tool.Bv2av(id)
	return RecordViewInfo(ctx, avid)
}
