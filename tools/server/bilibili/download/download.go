package download

import (
	"fmt"
	gcrawler "github.com/actliboy/hoper/server/go/lib/utils/generics/net/http/client/crawler"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client/crawler"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"tools/bilibili/api"
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

func FavReqs(page int) []*crawler.Request {
	var requests []*crawler.Request
	for i := 1; i <= page; i++ {
		req := gcrawler.NewRequest(api.GetFavListUrl(i), FavList)
		requests = append(requests, req)
	}
	return requests
}

func FavReq(page int) *crawler.Request {
	return gcrawler.NewRequest(api.GetFavListUrl(page), FavList)
}

var apiService = &api.API{}

func FavList(url string) ([]*crawler.Request, error) {
	res, err := api.Get[*api.FavList](url)
	if err != nil {
		return nil, err
	}
	var requests []*crawler.Request
	for _, fav := range res.Medias {
		aid := tool.Bv2av(fav.Bvid)
		req := GetRequestByFav(aid)
		requests = append(requests, req)
	}
	return requests, nil
}

func GetRequestByFav(aid int) *crawler.Request {
	return gcrawler.NewRequest(api.GetViewUrl(aid), ViewInfoHandleFun)
}

func ViewInfoHandleFun(url string) ([]*crawler.Request, error) {
	res, err := api.Get[api.ViewInfo](url)
	if err != nil {
		return nil, err
	}
	var requests []*crawler.Request
	for _, page := range res.Pages {
		video := &Video{res.Title, res.Aid, page.Cid, page.Page, page.Part, ""}
		req := gcrawler.NewRequest(api.GetPlayerUrl(res.Aid, page.Cid, 120), video.DownloadHandleFun)
		requests = append(requests, req)
	}
	return requests, nil
}

func (video *Video) DownloadHandleFun(url string) ([]*crawler.Request, error) {
	res, err := api.Get[*api.VideoInfo](url)
	if err != nil {
		return nil, err
	}
	if res.Quality != 120 {
		return []*crawler.Request{gcrawler.NewRequest(api.GetPlayerUrl(video.Aid, video.Cid, res.AcceptQuality[0]), video.DownloadHandleFun)}, nil
	}
	video.Quality = res.AcceptDescription[0]
	var requests []*crawler.Request
	for _, durl := range res.Durl {
		req := gcrawler.NewRequest(durl.Url, video.GetDownloadHandleFun(durl.Order))
		requests = append(requests, req)
	}
	return requests, nil
}

var _startUrlTem = "https://api.bilibili.com/x/web-interface/view?aid=%d"

func (video *Video) GetDownloadHandleFun(order int) crawler.HandleFun {
	referer := fmt.Sprintf(_startUrlTem, video.Aid)
	for i := 1; i <= video.Page; i++ {
		referer += fmt.Sprintf("/?p=%d", i)
	}

	return func(url string) ([]*crawler.Request, error) {

		c := http.Client{CheckRedirect: genCheckRedirectfun(referer)}

		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatalln(url, err)
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
		request.Header.Set("Cookie", api.Cookie)

		resp, err := c.Do(request)
		if err != nil {
			log.Fatalf("下载 %d 时出错, 错误信息：%s", video.Cid, err)
			return nil, err
		}

		if resp.StatusCode != http.StatusPartialContent {
			log.Fatalf("下载 %d 时出错, 错误码：%d", video.Cid, resp.StatusCode)
			return nil, fmt.Errorf("错误码： %d", resp.StatusCode)
		}
		defer resp.Body.Close()

		aidPath := tool.GetAidFileDownloadDir()
		filename := fmt.Sprintf("%d_%s_%d_%d.flv", video.Aid, video.Title, video.Page, order)
		file, err := os.Create(filepath.Join(aidPath, filename))
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
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
	return gcrawler.NewRequest(api.GetUpSpaceListUrl(upid, 1), UpSpaceListFirstPageHandleFun(upid))
}

func UpSpaceListFirstPageHandleFun(upid int) crawler.HandleFun {
	return func(url string) ([]*crawler.Request, error) {
		res, err := api.Get[*api.UpSpaceList](url)
		if err != nil {
			return nil, err
		}
		var requests []*crawler.Request
		for i := 1; i <= res.Page.Count; i++ {
			requests = append(requests, gcrawler.NewRequest(api.GetUpSpaceListUrl(upid, i), UpSpaceListHandleFun))
		}
		return requests, nil
	}
}

func UpSpaceListHandleFun(url string) ([]*crawler.Request, error) {
	res, err := api.Get[*api.UpSpaceList](url)
	if err != nil {
		return nil, err
	}
	var requests []*crawler.Request
	for _, video := range res.List.Vlist {
		req := GetRequestByFav(video.Aid)
		requests = append(requests, req)
	}
	return requests, nil
}
