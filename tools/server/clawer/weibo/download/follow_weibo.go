package download

import (
	"context"
	"github.com/liov/hoper/server/go/lib/utils/generics/net/http/client/crawler"
	"github.com/liov/hoper/server/go/lib/utils/log"
	"strings"
	"tools/clawer/weibo/rpc"
)

func GetUserFollowWeiboReq(maxId string) *crawler.Request {
	return &crawler.Request{
		TaskMeta: crawler.TaskMeta{BaseTaskMeta: crawler.BaseTaskMeta{Key: "GetUserFollowWeiboReq " + maxId}, Kind: KindGet},
		TaskFunc: func(ctx context.Context) ([]*crawler.Request, error) {
			log.Info("GetUserFollowWeiboReq:", maxId)
			list, err := rpc.GetFollowsWeibo(maxId)
			if err != nil {
				if strings.HasPrefix(err.Error(), "status:403") {

				}
				return nil, err
			}

			var requests []*crawler.Request

			if list.Statuses != nil {
				for _, mblog := range list.Statuses {
					if mblog.PicNum > 9 || mblog.IsLongText {
						requests = append(requests, LongWeiboReq(mblog.Id, false))
					} else {
						requests = append(requests, GetWeiboReq(mblog, false)...)
					}

					if mblog.RetweetedStatus != nil {

					}
				}
			}
			/*			if list != nil {
						requests = append(requests, GetUserFollowWeiboReq(list.MaxIdStr))
					}*/
			return requests, nil
		},
	}
}
