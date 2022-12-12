package download

import (
	"context"
	"github.com/liov/hoper/server/go/lib/utils/generics/conctrl"
	"github.com/liov/hoper/server/go/lib/utils/generics/net/http/client/crawler"
	"tools/clawer/weibo/rpc"
)

func LongWeiboReq(wid string, record bool) *crawler.Request {
	return &crawler.Request{
		TaskMeta: crawler.TaskMeta{
			BaseTaskMeta:   crawler.BaseTaskMeta{},
			Kind:           KindGet,
			TaskStatistics: conctrl.TaskStatistics{},
		},
		TaskFunc: func(ctx context.Context) ([]*crawler.Request, error) {
			mblog, err := rpc.GetLongWeibo(wid)
			if err != nil {
				return nil, err
			}
			return GetWeiboReq(mblog, record), nil
		},
	}
}
