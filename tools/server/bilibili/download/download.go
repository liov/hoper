package download

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/actliboy/hoper/server/go/lib/utils/dao/db/postgres"
	"github.com/actliboy/hoper/server/go/lib/utils/fs"
	gcrawler "github.com/actliboy/hoper/server/go/lib/utils/generics/net/http/client/crawler"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client/crawler"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
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
	Quality string
}

func FavReqs(pageStart, pageEnd int) []*crawler.Request {
	var requests []*crawler.Request
	for i := pageStart; i <= pageEnd; i++ {
		req := gcrawler.NewRequest(rpc.GetFavListUrl(i), FavList)
		requests = append(requests, req)
	}
	return requests
}

var apiService = &rpc.API{}

func FavList(ctx context.Context, url string) ([]*crawler.Request, error) {
	res, err := rpc.Get[*rpc.FavList](url)
	if err != nil {
		return nil, err
	}
	var requests []*crawler.Request
	for _, fav := range res.Medias {
		aid := tool.Bv2av(fav.Bvid)
		req1 := GetViewInfoReq(aid).SetKind(1)
		req2 := crawler.NewRequest(fav.Cover, DownloadCover(fav.Id)).SetKind(2)
		requests = append(requests, req1, req2)
	}
	return requests, nil
}

func GetViewInfoReq(aid int) *crawler.Request {
	return gcrawler.NewRequest(rpc.GetViewUrl(aid), ViewInfoHandleFun)
}

func ViewInfoHandleFun(ctx context.Context, url string) ([]*crawler.Request, error) {
	res, err := rpc.Get[rpc.ViewInfo](url)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	bilibiliDao := dao.NewDao(ctx, dao.Dao.Hoper.DB)
	err = bilibiliDao.CreateView(&dao.View{
		Bvid:        res.Bvid,
		Aid:         res.Aid,
		Data:        data,
		CoverRecord: false,
	})
	if err != nil && !postgres.IsDuplicate(err) {
		return nil, err
	}
	var requests []*crawler.Request
	for _, page := range res.Pages {
		video := &Video{fs.PathClean(res.Title), res.Aid, page.Cid, page.Page, page.Part, ""}
		_, err = video.DownloadHandleFun(ctx, rpc.GetPlayerUrl(res.Aid, page.Cid, 120))
		if err != nil {
			return nil, err
		}
		/*		req := crawler.NewRequest(rpc.GetPlayerUrl(res.Aid, page.Cid, 120), video.DownloadHandleFun)
				requests = append(requests, req)*/
	}
	return requests, nil
}

func (video *Video) DownloadHandleFun(ctx context.Context, url string) ([]*crawler.Request, error) {
	res, err := rpc.Get[*rpc.VideoInfo](url)
	if err != nil {
		return nil, err
	}

	video.Quality = res.AcceptDescription[0]
	var requests []*crawler.Request
	for _, durl := range res.Durl {
		req := gcrawler.NewRequest(durl.Url, video.GetDownloadHandleFun(durl.Order)).SetKind(3)
		requests = append(requests, req)
	}

	res.JsonClean()
	data, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	bilibiliDao := dao.NewDao(ctx, dao.Dao.Hoper.DB)
	err = bilibiliDao.CreateVideo(&dao.Video{
		Aid:    video.Aid,
		Cid:    video.Cid,
		Data:   data,
		Record: false,
	})
	if err != nil && !postgres.IsDuplicate(err) {
		return nil, err
	}

	return requests, nil
}

func (video *Video) GetDownloadHandleFun(order int) crawler.HandleFun {
	referer := rpc.GetViewUrl(video.Aid)
	for i := 1; i <= video.Page; i++ {
		referer += fmt.Sprintf("/?p=%d", i)
	}

	return func(ctx context.Context, url string) ([]*crawler.Request, error) {

		c := http.Client{CheckRedirect: genCheckRedirectfun(referer)}

		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Printf(url, err)
			return nil, err
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
			return nil, err
		}

		if resp.StatusCode != http.StatusPartialContent {
			log.Printf("下载 %d 时出错, 错误码：%d", video.Cid, resp.StatusCode)
			return nil, fmt.Errorf("错误码： %d", resp.StatusCode)
		}
		defer resp.Body.Close()

		filename := fmt.Sprintf("%d_%s_%s_%d_%d.flv", video.Aid, video.Title, video.Quality, video.Page, order)
		filename = strings.ReplaceAll(filename, " ", "")
		file, err := os.Create(filepath.Join(config.Conf.Bilibili.DownloadPath, filename))
		if err != nil {
			log.Println("错误信息：", err)
			return nil, err
		}
		defer file.Close()

		log.Println("正在下载："+filename, "质量：", video.Quality)
		_, err = io.Copy(file, resp.Body)
		if err != nil {
			log.Printf("下载失败 aid: %d, cid: %d, title: %s, part: %s",
				video.Aid, video.Cid, video.Title, video.Part)
			log.Println("错误信息：", err)

			// request again
			//go requestLater(file, resp, video)
			return nil, err
		}
		log.Println("下载完成：" + filename)

		return nil, nil
	}
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

func UpSpaceList(upid int) *crawler.Request {
	return gcrawler.NewRequest(rpc.GetUpSpaceListUrl(upid, 1), UpSpaceListFirstPageHandleFun(upid))
}

func UpSpaceListFirstPageHandleFun(upid int) crawler.HandleFun {
	return func(ctx context.Context, url string) ([]*crawler.Request, error) {
		res, err := rpc.Get[*rpc.UpSpaceList](url)
		if err != nil {
			return nil, err
		}
		var requests []*crawler.Request
		for i := 1; i <= res.Page.Count; i++ {
			requests = append(requests, gcrawler.NewRequest(rpc.GetUpSpaceListUrl(upid, i), UpSpaceListHandleFun))
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
		req := GetViewInfoReq(video.Aid)
		requests = append(requests, req)
	}
	return requests, nil
}

func GetByBvId(id string) *crawler.Request {
	avid := tool.Bv2av(id)
	return GetViewInfoReq(avid)
}

func DownloadCover(id int) crawler.HandleFun {
	return func(ctx context.Context, url string) ([]*crawler.Request, error) {
		return nil, client.DownloadImage(filepath.Join(config.Conf.Bilibili.DownloadPicPath, strconv.Itoa(id)+"_"+path.Base(url)), url)
	}
}
