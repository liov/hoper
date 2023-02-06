package download

import (
	"context"
	"github.com/liov/hoper/server/go/lib/utils/log"
	stringsi "github.com/liov/hoper/server/go/lib/utils/strings"
	"github.com/liov/hoper/server/go/lib_v2/utils/net/http/client/crawler"
	"strconv"
	"strings"
	"time"
	claweri "tools/clawer"
	"tools/clawer/weibo/config"
	"tools/clawer/weibo/rpc"
)

// 可下原图 Deprecated
func DownloadUserAllPhotoReq(uid int) *crawler.Request {
	return &crawler.Request{
		TaskMeta: crawler.TaskMeta{BaseTaskMeta: crawler.BaseTaskMeta{Key: strconv.Itoa(uid)}, Kind: KindNormal},
		TaskFunc: func(ctx context.Context) ([]*crawler.Request, error) {
			return []*crawler.Request{DownloadUserPhotoReq(uid, 1)}, nil
		},
	}
}

func DownloadUserPhotoReq(uid, page int) *crawler.Request {
	return &crawler.Request{
		TaskMeta: crawler.TaskMeta{BaseTaskMeta: crawler.BaseTaskMeta{Key: strconv.Itoa(uid) + " " + strconv.Itoa(page) + "DownloadUserPhotoReq"}, Kind: KindGet},
		TaskFunc: func(ctx context.Context) ([]*crawler.Request, error) {
			log.Infof("DownloadUserPhotoReq %d 第%d页", uid, page)
			piccards, err := rpc.GetPhotos(uid, page)
			if err != nil {
				if strings.HasPrefix(err.Error(), "status:403") {

				}
				return nil, err
			}
			var requests []*crawler.Request
			if piccards.Cards != nil {
				requests = DownloadPhotosReqs(uid, piccards.Cards)
			}

			if len(requests) != 0 {
				req := DownloadUserPhotoReq(uid, page+1)
				req.SetPriority(-1)
				requests = append(requests, req)
			}
			return requests, nil
		},
	}
}

func DownloadPhotosReqs(uid int, cards []*rpc.PicCardGroup) []*crawler.Request {
	var requests []*crawler.Request
	for _, card := range cards {
		requests = append(requests, DownloadPhotoReqs(uid, card)...)
	}
	return requests
}

func DownloadPhotoReqs(uid int, card *rpc.PicCardGroup) []*crawler.Request {
	var requests []*crawler.Request
	created := time.Now()
	for _, pic := range card.Pics {
		if pic.Type == "livephoto" {
			requests = append(requests, DownloadPhotoReq(created, uid, pic.Mblog.Id, pic.Video))
		}
		var url string
		if pic.PicBig != "" {
			url = pic.PicBig
		} else if pic.PicMiddle != "" {
			url = pic.PicBig
		}

		requests = append(requests, DownloadPhotoReq(created, uid, pic.Mblog.Id, url))
	}

	return requests
}

func DownloadPhotoReq(created time.Time, uid int, wid, url string) *crawler.Request {
	return &crawler.Request{
		TaskMeta: crawler.TaskMeta{BaseTaskMeta: crawler.BaseTaskMeta{Key: url}, Kind: KindDownload},
		TaskFunc: func(ctx context.Context) ([]*crawler.Request, error) {

			return nil, DownloadPhoto(created, uid, wid, url)
		},
	}
}

func DownloadPhoto(created time.Time, uid int, wid, url string) error {
	var baseUrl string
	typ := 1
	if strings.HasSuffix(url, "mov") {
		typ = 2
		baseUrl = stringsi.CountdownCutoff(url, "%2F")
	} else {
		baseUrl = stringsi.CountdownCutoff(url, "/")
		baseUrl = stringsi.Cutoff(baseUrl, "?")
	}

	return (&claweri.DownloadMeta{
		Dir: claweri.Dir{
			Platform: 4,
			PubAt:    created,
			Type:     typ,
			UserId:   uid,
			KeyIdStr: wid,
			BaseUrl:  baseUrl,
		},
		DownloadPath: config.Conf.Weibo.DownloadPath,
		Url:          url,
		Referer:      Referer,
	}).Download()

}
