package download

import (
	"context"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client/crawler"
	"tools/bilibili/rpc"
	"tools/bilibili/tool"
)

func CoverViewInfoHandleFun(ctx context.Context, url string) ([]*crawler.Request, error) {
	res, err := rpc.Get[rpc.ViewInfo](url)
	if err != nil && err.Error() != rpc.ErrorNotFound && err.Error() != rpc.ErrorNotPermission {
		return nil, err
	}

	return []*crawler.Request{crawler.NewUrlKindRequest(res.Pic, KindDownloadCover, CoverDownload(ctx, res.Owner.Mid, res.Aid))}, nil
}

func CoverFavList(ctx context.Context, url string) ([]*crawler.Request, error) {
	res, err := rpc.Get[*rpc.FavResourceList](url)
	if err != nil {
		return nil, err
	}
	var requests []*crawler.Request
	for _, fav := range res.Medias {
		aid := tool.Bv2av(fav.Bvid)
		req1 := GetViewInfoReq(aid, ViewInfoHandleFun)
		req2 := crawler.NewUrlKindRequest(fav.Cover, KindDownloadCover, CoverDownload(ctx, fav.Upper.Mid, fav.Id))
		requests = append(requests, req1, req2)
	}
	return requests, nil
}
