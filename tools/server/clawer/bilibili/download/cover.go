package download

import (
	"context"
	"github.com/liov/hoper/server/go/lib_v2/utils/net/http/client/crawler"
	"gorm.io/gorm"
	"time"
	claweri "tools/clawer"

	"path"
	"strings"
	"tools/clawer/bilibili/config"
	"tools/clawer/bilibili/dao"
	"tools/clawer/bilibili/rpc"
	"tools/clawer/bilibili/tool"
)

// 单个视频封面下载
func CoverViewInfoHandleFun(ctx context.Context, pubAt int, url string) ([]*crawler.Request, error) {
	res, err := rpc.Get[rpc.ViewInfo](url)
	if err != nil && err.Error() != rpc.ErrorNotFound && err.Error() != rpc.ErrorNotPermission {
		return nil, err
	}

	return []*crawler.Request{CoverDownloadReq(pubAt, res.Pic, res.Owner.Mid, res.Aid)}, nil
}

// 收藏夹封面下载
func DownloadFavCover(ctx context.Context, url string) ([]*crawler.Request, error) {
	res, err := rpc.Get[*rpc.FavResourceList](url)
	if err != nil {
		return nil, err
	}
	var requests []*crawler.Request
	for _, fav := range res.Medias {
		aid := tool.Bv2av(fav.Bvid)
		bilibiliDao := dao.NewDao(ctx, dao.Dao.Hoper.DB)
		view, err := bilibiliDao.ViewInfo(aid)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
		if err == gorm.ErrRecordNotFound {
			if !strings.HasSuffix(fav.Cover, "be27fd62c99036dce67efface486fb0a88ffed06.jpg") {
				req := CoverDownloadReq(fav.Pubtime, fav.Cover, fav.Upper.Mid, fav.Id)
				requests = append(requests, req)
			}
		} else {
			if view.Record == 0 {
				req := CoverDownloadReq(fav.Pubtime, fav.Cover, fav.Upper.Mid, fav.Id)
				requests = append(requests, req)
			}
		}
	}
	return requests, nil
}

func CoverDownloadReq(pubAt int, url string, upId, id int) *crawler.Request {
	return crawler.NewUrlKindRequest(url, KindDownloadCover, func(ctx context.Context, url string) ([]*crawler.Request, error) {
		return nil, CoverDownload(ctx, pubAt, url, upId, id)
	})
}

const NULLCOVER = "be27fd62c99036dce67efface486fb0a88ffed06.jpg"

func CoverDownload(ctx context.Context, pubAt int, url string, upId, id int) error {
	if strings.HasSuffix(url, NULLCOVER) {
		return nil
	}
	/*	var record bool
		dao.Dao.Hoper.Table(dao.TableNameView).Select("cover_record").Where("aid = ?", id).Scan(&record)
		if record {
			return nil
		}*/
	pubTime := time.Unix(int64(pubAt), 0)
	meta := claweri.DownloadMeta{
		Dir: claweri.Dir{
			Platform: 3,
			UserId:   upId,
			KeyId:    id,
			BaseUrl:  path.Base(url),
			Type:     1,
			PubAt:    pubTime,
		},
		DownloadPath: config.Conf.Bilibili.DownloadPicPath,
		Url:          url,
		Referer:      "",
	}

	return meta.Download(dao.Dao.Hoper.DB)
}
