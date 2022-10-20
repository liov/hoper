package download

import (
	"context"
	"encoding/json"
	"github.com/actliboy/hoper/server/go/lib/utils/conctrl"
	"github.com/actliboy/hoper/server/go/lib/utils/dao/db/postgres"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client/crawler"
	"gorm.io/gorm"
	"time"
	"tools/bilibili/dao"
	"tools/bilibili/rpc"
)

func (video *Video) RecordVideoReqAfterDownloadVideo() *crawler.Request {
	return &crawler.Request{
		TaskMeta: conctrl.TaskMeta{Kind: KindGetPlayerUrl},
		Key:      "",
		TaskFunc: func(ctx context.Context) ([]*crawler.Request, error) {
			videoInfo, err := video.RecordVideo(ctx)
			if err != nil {
				return nil, err
			}
			return GetDownloadRequests(videoInfo, video)
		},
	}
}

func (video *Video) GetVideoReqAfterDownloadVideo() *crawler.Request {
	return &crawler.Request{
		TaskMeta: conctrl.TaskMeta{Kind: KindGetPlayerUrl},
		Key:      "",
		TaskFunc: func(ctx context.Context) ([]*crawler.Request, error) {
			res, err := apiservice.GetPlayerInfo(video.Aid, video.Cid)
			if err != nil {
				if err.Error() == rpc.ErrorNotFound {
					dao.Dao.Hoper.Table(dao.TableNameVideo).Where(`cid = ?`, video.Cid).UpdateColumn("deleted_at", time.Now())
					return nil, nil
				}
				return nil, err
			}

			video.Quality = res.Quality

			return GetDownloadRequests(res, video)
		},
	}
}

func GetDownloadRequests(videoInfo *rpc.VideoInfo, video *Video) ([]*crawler.Request, error) {
	var requests []*crawler.Request
	for _, durl := range videoInfo.Durl {
		req := video.DownloadVideoReq("", VideoTypeFlv, durl.Order, durl.Url)
		requests = append(requests, req)
	}
	if videoInfo.Dash != nil {
		for _, v := range videoInfo.Dash.Video {
			if v.Id == video.Quality && v.Codecid == 12 {
				req := video.DownloadVideoReq("video", VideoTypeM4s, 1, v.BaseUrl)
				requests = append(requests, req)
				break
			}
		}
		req := video.DownloadVideoReq("audio", VideoTypeM4s, 1, videoInfo.Dash.Audio[0].BaseUrl)
		requests = append(requests, req)
	}
	return requests, nil
}

func (video *Video) RecordVideoReq() *crawler.Request {
	return &crawler.Request{
		TaskMeta: conctrl.TaskMeta{Kind: KindGetPlayerUrl},
		Key:      "",
		TaskFunc: func(ctx context.Context) ([]*crawler.Request, error) {
			_, err := video.RecordVideo(ctx)
			return nil, err
		},
	}
}

func (video *Video) RecordVideo(ctx context.Context) (*rpc.VideoInfo, error) {
	var dvideo dao.Video
	err := dao.Dao.Hoper.Table(dao.TableNameVideo).Select("cid,record").Where("cid = ?", video.Cid).First(&dvideo).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if dvideo.Cid > 0 {
		return nil, nil
	}
	res, err := apiservice.GetPlayerInfo(video.Aid, video.Cid)
	if err != nil {
		if err.Error() == rpc.ErrorNotFound {
			dao.Dao.Hoper.Table(dao.TableNameVideo).Where(`cid = ?`, video.Cid).UpdateColumn("deleted_at", time.Now())
			return nil, nil
		}
		return nil, err
	}

	video.Quality = res.Quality
	var durl []*rpc.Durl
	var dash *rpc.Dash
	if dvideo.Record == 0 {
		durl = res.Durl
		dash = res.Dash
	}

	bilibiliDao := dao.NewDao(ctx, dao.Dao.Hoper.DB)
	if err == gorm.ErrRecordNotFound || dvideo.Cid == 0 {
		res.JsonClean()
		data, err := json.Marshal(res)
		if err != nil {
			return nil, err
		}
		err = bilibiliDao.CreateVideo(&dao.Video{
			Aid:  video.Aid,
			Cid:  video.Cid,
			Data: data,
		})
		if err != nil && !postgres.IsDuplicate(err) {
			return nil, err
		}
	}

	res.Durl = durl
	res.Dash = dash

	return res, nil
}

func DownloadRecordVideo(engine *crawler.Engine) {
	now := time.Now()
	for {
		var videos []*Video
		dao.Dao.Hoper.DB.Raw(`SELECT a.owner->'mid' up_id,b.aid,b.cid,a.title,a.p->'page' page,a.p->'part' part
FROM `+dao.TableNameVideo+` b 
LEFT JOIN (SELECT data->'title' title, data->'owner' owner, jsonb_path_query(data,'$.pages[*]') p FROM `+dao.TableNameView+`)  a ON (a.p->'cid')::int8 = b.cid
WHERE b.record = false AND b.created_at < ? AND b.`+postgres.NotDeleted+` ORDER BY b.created_at DESC LIMIT 20`, now).Find(&videos)
		if len(videos) == 0 {
			return
		}
		for _, video := range videos {
			if video.Title == "" {
				req := ViewRecordUpdateReqAfterRecordVideo(video.Aid)
				engine.Engine.AddTask(engine.NewTask(req))
			} else {
				req := video.GetVideoReqAfterDownloadVideo()
				engine.Engine.AddTask(engine.NewTask(req))
			}
		}
	}
}
