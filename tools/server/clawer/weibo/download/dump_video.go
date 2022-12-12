package download

import (
	"context"
	"github.com/liov/hoper/server/go/lib/utils/generics/net/http/client/crawler"
	"github.com/liov/hoper/server/go/lib/utils/log"
	"github.com/liov/hoper/server/go/lib/utils/net/http/client"
	stringsi "github.com/liov/hoper/server/go/lib/utils/strings"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"tools/clawer/weibo/config"
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
			return DownloadVideoWarpReq(int(mblog.User.Id), mblog.Id, url)
		}
	}
	return nil
}

func DownloadVideoWarpReq(uid int, wid, url string) *crawler.Request {
	return &crawler.Request{
		TaskMeta: crawler.TaskMeta{BaseTaskMeta: crawler.BaseTaskMeta{Key: url}, Kind: KindDownload},
		TaskFunc: func(ctx context.Context) ([]*crawler.Request, error) {
			return nil, DownloadVideo(uid, wid, url)
		},
	}
}

func DownloadVideo(uid int, wid, url string) error {
	baseUrl := stringsi.CountdownCutoff(stringsi.CutoffContain(url, "mp4"), "/")
	filepath := filepath.Join(config.Conf.Weibo.DownloadVideoPath, strconv.Itoa(uid), strconv.Itoa(uid)+"_"+wid+"_"+baseUrl)
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		err = client.DownloadFileWithRefer(filepath, url, Referer)
		if err != nil {
			log.Info("下载视频失败：", err)
			return err
		}
		log.Info("下载视频成功：", filepath)
	} else {
		log.Info("视频已存在：", filepath)
	}
	return nil
}
