package download

import (
	"context"
	"github.com/liov/hoper/server/go/lib/utils/generics/net/http/client/crawler"
	"gorm.io/gorm"

	"github.com/liov/hoper/server/go/lib/utils/net/http/client"

	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"tools/clawer/bilibili/config"
	"tools/clawer/bilibili/dao"
	"tools/clawer/bilibili/rpc"
	"tools/clawer/bilibili/tool"
)

// 单个视频封面下载
func CoverViewInfoHandleFun(ctx context.Context, url string) ([]*crawler.Request, error) {
	res, err := rpc.Get[rpc.ViewInfo](url)
	if err != nil && err.Error() != rpc.ErrorNotFound && err.Error() != rpc.ErrorNotPermission {
		return nil, err
	}

	return []*crawler.Request{CoverDownloadReq(res.Pic, res.Owner.Mid, res.Aid)}, nil
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
				req := CoverDownloadReq(fav.Cover, fav.Upper.Mid, fav.Id)
				requests = append(requests, req)
			}
		} else {
			if view.CoverRecord == false {
				req := CoverDownloadReq(fav.Cover, fav.Upper.Mid, fav.Id)
				requests = append(requests, req)
			}
		}
	}
	return requests, nil
}

func CoverDownloadReq(url string, upId, id int) *crawler.Request {
	return crawler.NewUrlKindRequest(url, KindDownloadCover, func(ctx context.Context, url string) ([]*crawler.Request, error) {
		return nil, CoverDownload(ctx, url, upId, id)
	})
}

const NULLCOVER = "be27fd62c99036dce67efface486fb0a88ffed06.jpg"

func CoverDownload(ctx context.Context, url string, upId, id int) error {
	if strings.HasSuffix(url, NULLCOVER) {
		return nil
	}
	/*	var record bool
		dao.Dao.Hoper.Table(dao.TableNameView).Select("cover_record").Where("aid = ?", id).Scan(&record)
		if record {
			return nil
		}*/
	filepath := filepath.Join(config.Conf.Bilibili.DownloadPicPath, strconv.Itoa(upId), strconv.Itoa(upId)+"_"+strconv.Itoa(id)+"_"+path.Base(url))
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		err = client.DownloadImage(filepath, url)
		if err != nil {
			log.Println("下载图片失败：", err)
			return err
		}
		dao.Dao.Hoper.Table(dao.TableNameView).Where("aid = ?", id).Update("cover_record", true)
		log.Println("下载图片成功：", filepath)
	} else {
		log.Println("图片已存在：", filepath)
	}
	return nil
}
