package download

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/actliboy/hoper/server/go/lib/utils/conctrl"
	"github.com/actliboy/hoper/server/go/lib/utils/dao/db/postgres"
	"github.com/actliboy/hoper/server/go/lib/utils/fs"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client/crawler"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"
	"tools/bilibili/config"
	"tools/bilibili/dao"
	"tools/bilibili/rpc"
	"tools/bilibili/tool"
)

type Video struct {
	Title   string
	Aid     int
	Cid     int
	Page    int
	Part    string
	Quality int
}

const (
	KindGetFavListUrl conctrl.Kind = 0
	KindViewInfo      conctrl.Kind = 1
	KindDownloadCover conctrl.Kind = 2
	KindGetPlayerUrl  conctrl.Kind = 3
	KindDownloadVideo conctrl.Kind = 4
)

func FavList(ctx context.Context, url string) ([]*crawler.Request, error) {
	res, err := rpc.Get[*rpc.FavResourceList](url)
	if err != nil {
		return nil, err
	}
	var requests []*crawler.Request
	for _, fav := range res.Medias {
		aid := tool.Bv2av(fav.Bvid)
		req1 := GetViewInfoReq(aid, ViewInfoHandleFun)
		req2 := crawler.NewUrlKindRequest(fav.Cover, KindDownloadCover, CoverDownload(ctx, fav.Id))
		requests = append(requests, req1, req2)
	}
	return requests, nil
}

func ViewInfoHandleFun(ctx context.Context, url string) ([]*crawler.Request, error) {
	res, err := rpc.Get[rpc.ViewInfo](url)
	if err != nil && err.Error() != rpc.ErrorNotFound && err.Error() != rpc.ErrorNotPermission {
		return nil, err
	}
	bilibiliDao := dao.NewDao(ctx, dao.Dao.Hoper.DB)
	exists, err := bilibiliDao.ViewExists(res.Aid)
	if err != nil {
		return nil, err
	}
	if !exists {
		data, err := json.Marshal(res)
		if err != nil {
			return nil, err
		}
		err = bilibiliDao.CreateView(&dao.View{
			Bvid:        res.Bvid,
			Aid:         res.Aid,
			Data:        data,
			CoverRecord: false,
		})
		if err != nil && !postgres.IsDuplicate(err) {
			return nil, err
		}
	}
	var requests []*crawler.Request
	for _, page := range res.Pages {
		video := &Video{fs.PathClean(res.Title), res.Aid, page.Cid, page.Page, page.Part, 0}

		req := crawler.NewUrlKindRequest(rpc.GetPlayerUrl(res.Aid, page.Cid, 120), KindGetPlayerUrl, video.PlayerUrlHandleFun)
		requests = append(requests, req)
	}
	return requests, nil
}

func (video *Video) PlayerUrlHandleFun(ctx context.Context, url string) ([]*crawler.Request, error) {
	var dvideo dao.Video
	err := dao.Dao.Hoper.Table(dao.TableNameVideo).Select("cid,record").Where("cid = ?", video.Cid).First(&dvideo).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if dvideo.Record {
		return nil, nil
	}
	res, err := rpc.Get[*rpc.VideoInfo](url)
	if err != nil {
		if err.Error() == rpc.ErrorNotFound {
			dao.Dao.Hoper.Table(dao.TableNameVideo).Where(`cid = ?`, video.Cid).UpdateColumn("deleted_at", time.Now())
			return nil, nil
		}
		return nil, err
	}

	video.Quality = res.Quality
	if !dvideo.Record {
		for _, durl := range res.Durl {
			err = video.DownloadVideoHandleFun(durl.Order, durl.Url)
			if err != nil {
				return nil, err
			}
		}
	}

	bilibiliDao := dao.NewDao(ctx, dao.Dao.Hoper.DB)
	if err == gorm.ErrRecordNotFound || dvideo.Cid == 0 {
		res.JsonClean()
		data, err := json.Marshal(res)
		if err != nil {
			return nil, err
		}
		err = bilibiliDao.CreateVideo(&dao.Video{
			Aid:    video.Aid,
			Cid:    video.Cid,
			Data:   data,
			Record: false,
		})
		if err != nil && !postgres.IsDuplicate(err) {
			return nil, err
		}
	}

	return nil, nil
}

func (video *Video) DownloadVideoHandleFun(order int, url string) error {
	referer := rpc.GetViewUrl(video.Aid)
	referer = referer + fmt.Sprintf("/?p=%d", video.Page)

	c := http.Client{CheckRedirect: genCheckRedirectfun(referer)}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	request.Header.Set("User-Agent", client.UserAgent1)
	request.Header.Set("Accept", "*/*")
	request.Header.Set("Accept-Language", "en-US,en;q=0.5")
	request.Header.Set("Accept-Encoding", "gzip, deflate, br")
	request.Header.Set("Range", "bytes=0-")
	request.Header.Set("Referer", referer)
	request.Header.Set("Origin", "https://www.bilibili.com")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("Cookie", rpc.Cookie)

	resp, err := c.Do(request)
	if err != nil {
		log.Printf("下载 %d 时出错, 错误信息：%s", video.Cid, err)
		return err
	}

	if resp.StatusCode != http.StatusPartialContent {
		log.Printf("下载 %d 时出错, 错误码：%d", video.Cid, resp.StatusCode)
		return fmt.Errorf("错误码： %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	filename := fmt.Sprintf("%d_%d_%s_%d_%d.flv.downloading", video.Aid, video.Cid, video.Title, order, video.Quality)
	filename = fs.PathClean(filename)
	filename = filepath.Join(config.Conf.Bilibili.DownloadVideoPath, filename)
	file, err := os.Create(filename)
	if err != nil {
		log.Println("错误信息：", err)
		return err
	}

	newname := filename[:len(filename)-len(".downloading")]

	log.Println("正在下载："+filename, "质量：", video.Quality)
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		os.Remove(filename)
		log.Printf("下载失败 aid: %d, cid: %d, title: %s, part: %s",
			video.Aid, video.Cid, video.Title, video.Part)
		log.Println("错误信息：", err)

		// request again
		//go requestLater(file, resp, video)
		return err
	}
	file.Close()

	err = os.Rename(filename, newname)
	if err != nil {
		return err
	}
	dao.Dao.Hoper.Table(dao.TableNameVideo).Where("cid = ?", video.Cid).Update("record", true)
	log.Println("下载完成：" + newname)

	return nil
}

func genCheckRedirectfun(referer string) func(req *http.Request, via []*http.Request) error {
	return func(req *http.Request, via []*http.Request) error {
		req.Header.Set("Referer", referer)
		return nil
	}
}
func requestLater(file *os.File, resp *http.Response, video *Video) error {

	log.Println("连接失败，30秒后重试 (Unable to open the file due to the remote host, request in 30 seconds)")
	time.Sleep(time.Second * 30)

	_, err := io.Copy(file, resp.Body)
	if err != nil {
		log.Printf("下载失败 aid: %d, cid: %d, title: %s, part: %s again",
			video.Aid, video.Cid, video.Title, video.Part)
	}
	return err
}

func UpSpaceListFirstPageHandleFun(upid int) crawler.HandleFun {
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
		req := GetViewInfoReq(video.Aid, ViewInfoHandleFun)
		requests = append(requests, req)
	}
	return requests, nil
}

func CoverDownload(ctx context.Context, id int) crawler.HandleFun {
	var record bool
	dao.Dao.Hoper.Table(dao.TableNameView).Select("cover_record").Where("aid = ?", id).Scan(&record)
	if record {
		return nil
	}
	return func(ctx context.Context, url string) ([]*crawler.Request, error) {
		err := client.DownloadImage(filepath.Join(config.Conf.Bilibili.DownloadPicPath, strconv.Itoa(id)+"_"+path.Base(url)), url)
		if err != nil {
			log.Println("下载图片失败：", err)
			return nil, err
		}
		dao.Dao.Hoper.Table(dao.TableNameView).Where("aid = ?", id).Update("cover_record", true)
		return nil, nil
	}
}
