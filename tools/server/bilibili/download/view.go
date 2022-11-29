package download

import (
	"context"
	"encoding/json"
	gormpostgres "github.com/actliboy/hoper/server/go/lib/utils/dao/db/gorm/postgres"
	"github.com/actliboy/hoper/server/go/lib/utils/dao/db/postgres"

	"github.com/actliboy/hoper/server/go/lib/utils/generics/net/http/client/crawler"

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
			return ViewGetRecordVideoReqs(view)
		},
	}
}

func ViewGetRecordVideoReqs(view *rpc.ViewInfo) ([]*crawler.Request, error) {
	if view == nil {
		return nil, nil
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
	if err != nil {
		if err.Error() == rpc.ErrorNotFound || err.Error() == rpc.ErrorNotPermission {
			return nil, nil
		}
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
			view, err := ViewRecordUpdate(ctx, aid)
			if err != nil {
				return nil, err
			}
			return ViewGetRecordVideoReqs(view)
		},
	}
}

func ViewRecordUpdateReq(aid int) *crawler.Request {
	return &crawler.Request{
		TaskMeta: crawler.TaskMeta{Kind: KindViewInfo},
		TaskFunc: func(ctx context.Context) ([]*crawler.Request, error) {
			_, err := ViewRecordUpdate(ctx, aid)
			return nil, err
		},
	}
}

func ViewRecordUpdate(ctx context.Context, aid int) (*rpc.ViewInfo, error) {
	bilibiliDao := dao.NewDao(ctx, dao.Dao.Hoper.DB)
	exists, err := bilibiliDao.ViewExists(aid)
	if err != nil {
		return nil, err
	}
	if !exists {
		return RecordViewInfo(ctx, aid)
	}

	err = dao.Dao.Hoper.Exec(`INSERT INTO `+dao.TableNameViewBak+`(aid,data) (SELECT aid,data FROM `+dao.TableNameView+` WHERE aid = ?) `, aid).Error
	/*	if err != nil {
		return nil, err
	}*/
	res, err := apiservice.GetView(aid)
	if err != nil {
		if err.Error() == rpc.ErrorNotFound || err.Error() == rpc.ErrorNotPermission {
			gormpostgres.DeleteSQL(dao.Dao.Hoper.DB, dao.TableNameView, "aid", aid)
			return nil, nil
		}
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

	return res, nil
}

func RecordViewInfoByBvId(ctx context.Context, id string) (*rpc.ViewInfo, error) {
	avid := tool.Bv2av(id)
	return RecordViewInfo(ctx, avid)
}
