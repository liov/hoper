package download

import (
	"context"
	"github.com/liov/hoper/server/go/lib/utils/log"
	stringsi "github.com/liov/hoper/server/go/lib/utils/strings"
	timei "github.com/liov/hoper/server/go/lib/utils/time"
	"github.com/liov/hoper/server/go/lib_v2/utils/net/http/client/crawler"
	"strconv"
	"strings"
	"time"
	claweri "tools/clawer"
	"tools/clawer/weibo/config"
	"tools/clawer/weibo/dao"
	"tools/clawer/weibo/rpc"
)

func DownloadUserAllVideoReq(uid int) *crawler.Request {
	return &crawler.Request{
		TaskMeta: crawler.TaskMeta{BaseTaskMeta: crawler.BaseTaskMeta{Key: strconv.Itoa(uid)}, Kind: KindNormal},
		TaskFunc: func(ctx context.Context) ([]*crawler.Request, error) {
			return []*crawler.Request{DownloadUserVideoReq(uid, 1)}, nil
		},
	}
}

func DownloadUserVideoReq(uid, page int) *crawler.Request {
	return &crawler.Request{
		TaskMeta: crawler.TaskMeta{BaseTaskMeta: crawler.BaseTaskMeta{Key: strconv.Itoa(uid) + " " + strconv.Itoa(page) + "DownloadUserVideoReq"}, Kind: KindGet},
		TaskFunc: func(ctx context.Context) ([]*crawler.Request, error) {
			log.Infof("DownloadUserVideoReq %d 第%d页", uid, page)
			piccards, err := rpc.GetVideos(uid, page)
			if err != nil {
				if strings.HasPrefix(err.Error(), "status:403") {
					time.Sleep(time.Minute * 5)
				}
				return nil, err
			}
			var requests []*crawler.Request
			if piccards.Cards != nil {
				requests = DownloadVideosReq(piccards.Cards)
			}
			if len(requests) != 0 {
				req := DownloadUserVideoReq(uid, page+1)
				req.SetPriority(-1)
				requests = append(requests, req)
			}
			return requests, nil
		},
	}
}

func DownloadVideosReq(cards []*rpc.CardGroup) []*crawler.Request {
	var requests []*crawler.Request
	for _, card := range cards {
		if card.Mblog != nil {
			req := DownloadVideoReq(card.Mblog)
			if req != nil {
				requests = append(requests, req)
			}
		}
	}

	return requests
}

func DownloadVideoReq(mblog *rpc.Mblog) *crawler.Request {
	createdAt, _ := time.Parse(time.RubyDate, mblog.CreatedAt)
	created := createdAt.Format(timei.TimeFormatDisplay)

	if mblog.PageInfo != nil {
		var url string
		if mblog.PageInfo.Urls.Mp4720PMp4 != "" {
			url = mblog.PageInfo.Urls.Mp4720PMp4
		} else if mblog.PageInfo.Urls.Mp4HdMp4 != "" {
			url = mblog.PageInfo.Urls.Mp4HdMp4
		} else {
			url = mblog.PageInfo.Urls.Mp4LdMp4
		}
		if url != "" {
			return DownloadVideoWarpReq(created, mblog.User.Id, mblog.Id, url)
		}
	}
	return nil
}

func DownloadVideoWarpReq(created string, uid int, wid, url string) *crawler.Request {
	return &crawler.Request{
		TaskMeta: crawler.TaskMeta{BaseTaskMeta: crawler.BaseTaskMeta{Key: url}, Kind: KindDownload},
		TaskFunc: func(ctx context.Context) ([]*crawler.Request, error) {
			return nil, DownloadVideo(created, uid, wid, url)
		},
	}
}

func DownloadVideo(created string, uid int, wid, url string) error {
	baseUrl := stringsi.CountdownCutoff(stringsi.CutoffContain(url, "mp4"), "/")
	return (&claweri.DownloadMeta{
		Dir: claweri.Dir{
			PubAt:    created,
			Platform: 3,
			UserId:   uid,
			KeyIdStr: wid,
			BaseUrl:  baseUrl,
		},
		DownloadPath: config.Conf.Weibo.DownloadPath,
		Url:          url,
		Referer:      Referer,
	}).Download(dao.Dao.Hoper.DB)
}
