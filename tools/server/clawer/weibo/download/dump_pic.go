package download

import (
	"context"
	"github.com/liov/hoper/server/go/lib/utils/fs"
	"github.com/liov/hoper/server/go/lib/utils/log"
	"github.com/liov/hoper/server/go/lib/utils/net/http/client"
	stringsi "github.com/liov/hoper/server/go/lib/utils/strings"
	"github.com/liov/hoper/server/go/lib_v2/utils/net/http/client/crawler"
	"path/filepath"
	"strconv"
	"strings"
	"tools/clawer/weibo/config"
	"tools/clawer/weibo/rpc"
)

// 可下原图
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

	for _, pic := range card.Pics {
		if pic.Type == "livephoto" {
			requests = append(requests, DownloadPhotoReq(uid, pic.Mblog.Id, pic.Video))
		}
		var url string
		if pic.PicBig != "" {
			url = pic.PicBig
		} else if pic.PicMiddle != "" {
			url = pic.PicBig
		}
		requests = append(requests, DownloadPhotoReq(uid, pic.Mblog.Id, url))
	}

	return requests
}

func DownloadPhotoReq(uid int, wid, url string) *crawler.Request {
	return &crawler.Request{
		TaskMeta: crawler.TaskMeta{BaseTaskMeta: crawler.BaseTaskMeta{Key: url}, Kind: KindDownload},
		TaskFunc: func(ctx context.Context) ([]*crawler.Request, error) {

			return nil, DownloadPhoto(uid, wid, url)
		},
	}
}

func DownloadPhoto(uid int, wid, url string) error {
	var baseUrl string
	if strings.HasSuffix(url, "mov") {
		baseUrl = stringsi.CountdownCutoff(url, "%2F")
	} else {
		baseUrl = stringsi.CountdownCutoff(url, "/")
	}

	filepath := filepath.Join(config.Conf.Weibo.DownloadPicPath, strconv.Itoa(uid), strconv.Itoa(uid)+"_"+wid+"_"+baseUrl)

	if fs.NotExist(filepath) {
		err := client.DownloadFileWithRefer(filepath, url, Referer)
		if err != nil {
			log.Info("下载图片失败：", err)
			return err
		}
		log.Info("下载图片成功：", filepath)
	} else {
		log.Info("图片已存在：", filepath)
	}
	return nil
}
