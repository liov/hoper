package download

import (
	"context"
	"encoding/json"
	"github.com/liov/hoper/server/go/lib/utils/log"
	"github.com/liov/hoper/server/go/lib_v2/utils/net/http/client/crawler"
	"strconv"
	"strings"
	"time"
	"tools/clawer/weibo/rpc"
)

// 只能下缩略图
func DownloadUserPhotoReqV2(uid, page int) *crawler.Request {
	return &crawler.Request{
		TaskMeta: crawler.TaskMeta{BaseTaskMeta: crawler.BaseTaskMeta{Key: strconv.Itoa(uid) + " " + strconv.Itoa(page) + "DownloadUserPhotoReqV2"}, Kind: KindGet},
		TaskFunc: func(ctx context.Context) ([]*crawler.Request, error) {
			log.Infof("DownloadUserPhotoReqV2 %d 第%d页", uid, page)
			piccards, err := rpc.GetChannels(rpc.PIC, uid, page)
			if err != nil {
				if strings.HasPrefix(err.Error(), "status:403") {

				}
				if _, ok := err.(*json.UnmarshalTypeError); ok {

				}
				return nil, err
			}
			var requests []*crawler.Request
			if piccards.Cards != nil {
				requests = DownloadPhotosReqsV2(uid, piccards.Cards)
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

func DownloadPhotosReqsV2(uid int, cards []*rpc.CardGroup) []*crawler.Request {
	var requests []*crawler.Request
	for _, card := range cards {
		if card.Mblog != nil {
			requests = append(requests, DownloadPhotoReqsV2(card.Mblog)...)
		}
	}
	return requests
}

func DownloadPhotoReqsV2(mblog *rpc.Mblog) []*crawler.Request {
	var requests []*crawler.Request
	createdAt, _ := time.Parse(time.RubyDate, mblog.CreatedAt)
	for _, pic := range mblog.Pics {
		if pic.Type == "livephotos" {
			requests = append(requests, DownloadPhotoReq(createdAt, mblog.User.Id, mblog.Id, pic.VideoSrc))
		}
		var url string
		if pic.Large.Url != "" {
			url = pic.Large.Url
		}
		if url != "" {
			requests = append(requests, DownloadPhotoReq(createdAt, mblog.User.Id, mblog.Id, url))
		}
	}

	return requests
}
