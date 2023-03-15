package download

import (
	"context"
	"github.com/liov/hoper/server/go/lib/utils/net/http/client/crawler"
	"tools/clawer/bilibili/rpc"
)

func UpSpaceListFirstPageHandleFun(upid int) crawler.HandleFunc {
	return func(ctx context.Context, url string) ([]*crawler.Request, error) {
		res, err := rpc.Get[*rpc.UpSpaceList](url)
		if err != nil {
			return nil, err
		}
		var requests []*crawler.Request
		for i := 1; i <= res.Page.Count; i++ {
			requests = append(requests, crawler.NewUrlRequest(rpc.GetUpSpaceListUrl(upid, i), UpSpaceListHandleFun))
		}
		return requests, nil
	}
}

func UpSpaceListHandleFun(ctx context.Context, url string) ([]*crawler.Request, error) {
	res, err := rpc.Get[*rpc.UpSpaceList](url)
	if err != nil {
		return nil, err
	}
	var requests []*crawler.Request
	for _, video := range res.List.Vlist {
		req := RecordViewInfoReqAfterRecordVideo(video.Aid)
		requests = append(requests, req)
	}
	return requests, nil
}

func UpSpaceList(upid int, handleFun crawler.HandleFunc) *crawler.Request {
	return crawler.NewUrlRequest(rpc.GetUpSpaceListUrl(upid, 1), UpSpaceListFirstPageHandleFun(upid))
}
