package download

import (
	"context"
	"github.com/liov/hoper/server/go/lib/utils/generics/net/http/client/crawler"
	"github.com/liov/hoper/server/go/lib/utils/log"
	"github.com/liov/hoper/server/go/lib/utils/net/http/client"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"tools/clawer/weibo/config"
	"tools/clawer/weibo/rpc"
)

func GetUserAllFollowsReq(uid int) *crawler.Request {
	return GetUserFollowReq(uid, 1)
}

func GetUserFollowReq(uid, page int) *crawler.Request {
	return &crawler.Request{
		TaskMeta: crawler.TaskMeta{BaseTaskMeta: crawler.BaseTaskMeta{Key: strconv.Itoa(uid)}, Kind: KindGetPhoto},
		TaskFunc: func(ctx context.Context) ([]*crawler.Request, error) {
			log.Infof("GetUserFollowReq 第%d页", page)
			follow, err := rpc.GetFollows(config.Conf.Weibo.UserId, page)
			if err != nil {
				if strings.HasPrefix(err.Error(), "status:403") {

				}
				return nil, err
			}
			var requests []*crawler.Request
			for _, card := range follow.Cards {
				for _, group := range card.CardGroup {
					for _, user := range group.Elements {
						requests = append(requests, DownloadUserAllPhotoReq(user.Uid))
					}
				}
			}
			if len(requests) != 0 {
				requests = append(requests, DownloadUserPhotoReq(uid, page+1))
			}
			return requests, nil
		},
	}
}

func DownloadUserAllPhotoReq(uid int) *crawler.Request {
	return DownloadUserPhotoReq(uid, 1)
}

func DownloadUserPhotoReq(uid, page int) *crawler.Request {
	return &crawler.Request{
		TaskMeta: crawler.TaskMeta{BaseTaskMeta: crawler.BaseTaskMeta{Key: strconv.Itoa(uid)}, Kind: KindGetPhoto},
		TaskFunc: func(ctx context.Context) ([]*crawler.Request, error) {
			log.Infof("DownloadUserPhotoReq 第%d页", page)
			piccards, err := rpc.GetPhotos(uid, page)
			if err != nil {
				if strings.HasPrefix(err.Error(), "status:403") {

				}
				return nil, err
			}
			var requests []*crawler.Request
			if piccards.Cards != nil {
				for _, card := range piccards.Cards {
					requests = append(requests, DownloadPhotoReqs(uid, card)...)
				}
			}
			if len(requests) != 0 {
				requests = append(requests, DownloadUserPhotoReq(uid, page+1))
			}
			return requests, nil
		},
	}
}

func DownloadPhotoReqs(uid int, group *rpc.PicCardGroup) []*crawler.Request {
	var requests []*crawler.Request

	for _, pic := range group.Pics {
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
		TaskMeta: crawler.TaskMeta{BaseTaskMeta: crawler.BaseTaskMeta{Key: url}, Kind: KindDownloadPhoto},
		TaskFunc: func(ctx context.Context) ([]*crawler.Request, error) {

			return nil, DownloadPhoto(uid, wid, url)
		},
	}
}

func DownloadPhoto(uid int, wid, url string) error {
	filepath := filepath.Join(config.Conf.Weibo.DownloadPicPath, strconv.Itoa(uid), strconv.Itoa(uid)+"_"+wid+"_"+path.Base(url))
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		err = client.DownloadFileWithRefer(filepath, url, Referer)
		if err != nil {
			log.Info("下载图片失败：", err)
			return err
		}
		log.Info("下载图片成功：", filepath)
	}
	log.Info("图片已存在：", filepath)
	return nil
}