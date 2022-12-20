package download

import (
	"context"
	"github.com/liov/hoper/server/go/lib/utils/log"
	"github.com/liov/hoper/server/go/lib_v2/utils/net/http/client/crawler"
	"strconv"
	"strings"
	"time"
	"tools/clawer/weibo/config"
	"tools/clawer/weibo/rpc"
)

func GetUserAllFollowsReq(uid int) *crawler.Request {
	return &crawler.Request{
		TaskMeta: crawler.TaskMeta{BaseTaskMeta: crawler.BaseTaskMeta{Key: strconv.Itoa(uid) + "GetUserAllFollowsReq"}, Kind: KindNormal},
		TaskFunc: func(ctx context.Context) ([]*crawler.Request, error) {
			return []*crawler.Request{GetUserFollowReq(uid, 1)}, nil
		},
	}
}

func GetUserFollowReq(uid, page int) *crawler.Request {
	return &crawler.Request{
		TaskMeta: crawler.TaskMeta{BaseTaskMeta: crawler.BaseTaskMeta{Key: strconv.Itoa(uid) + " " + strconv.Itoa(page) + "GetUserFollowReq"}, Kind: KindGet},
		TaskFunc: func(ctx context.Context) ([]*crawler.Request, error) {
			log.Infof("GetUserFollowReq %d 第%d页", uid, page)
			follow, err := rpc.GetFollows(config.Conf.Weibo.UserId, page)
			if err != nil {
				if strings.HasPrefix(err.Error(), "status:403") {
					time.Sleep(time.Minute * 5)
				}
				return nil, err
			}
			var requests []*crawler.Request
			for _, card := range follow.Cards {
				for _, group := range card.CardGroup {
					for _, user := range group.Elements {
						//requests = append(requests, DownloadUserAllPhotoReq(user.Uid))
						requests = append(requests, DownloadUserAllVideoReq(user.Uid))
					}
				}
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
