package download

import (
	"context"
	"fmt"
	"github.com/actliboy/hoper/server/go/lib/utils/conctrl"
	"github.com/actliboy/hoper/server/go/lib/utils/fs"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client/crawler"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"tools/bilibili/config"
	"tools/bilibili/dao"
	"tools/bilibili/rpc"
)

type Video struct {
	UpId    int
	Title   string
	Aid     int
	Cid     int
	Page    int
	Part    string
	Quality int
}

func (video *Video) DownloadVideoReq(typ string, ext int, order int, url string) *crawler.Request {
	return &crawler.Request{
		TaskMeta: conctrl.TaskMeta{Kind: KindDownloadVideo},
		Key:      "",
		TaskFunc: func(ctx context.Context) ([]*crawler.Request, error) {
			return video.DownloadVideo(typ, ext, order, url)
		},
	}
}

func (video *Video) DownloadVideo(typ string, ext int, order int, url string) ([]*crawler.Request, error) {
	referer := rpc.GetViewUrl(video.Aid)
	referer = referer + fmt.Sprintf("/?p=%d", video.Page)

	c := http.Client{CheckRedirect: genCheckRedirectfun(referer)}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
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

	var filename string

	video.Part = fs.PathClean(video.Part)
	video.Title = fs.PathClean(video.Title)
	if video.Part == video.Title {
		video.Part = "!part=title!"
	}

	if ext == VideoTypeFlv {
		filename = fmt.Sprintf("%d_%d_%d_%s_%s_%d_%d.flv.downloading", video.UpId, video.Aid, video.Cid, video.Title, video.Part, order, video.Quality)
	} else {
		filename = fmt.Sprintf("%d_%d_%d_%s_%s_%d.m4s.%s.downloading", video.UpId, video.Aid, video.Cid, video.Title, video.Part, video.Quality, typ)
	}

	filename = filepath.Join(config.Conf.Bilibili.DownloadVideoPath, strconv.Itoa(video.UpId), filename)
	file, err := fs.Create(filename)
	if err != nil {
		log.Println("错误信息：", err)
		return nil, err
	}

	newname := filename[:len(filename)-len(".downloading")]

	log.Println("正在下载："+filename, "质量：", video.Quality)
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		os.Remove(filename)
		log.Printf("下载失败 filename: %s", filename)
		log.Println("错误信息：", err)

		// request again
		//go requestLater(file, resp, video)
		return nil, err
	}
	file.Close()

	err = os.Rename(filename, newname)
	if err != nil {
		return nil, err
	}
	if ext == VideoTypeFlv {
		dao.Dao.Hoper.Table(dao.TableNameVideo).Where("cid = ?", video.Cid).Update("record", 1)
	}

	log.Println("下载完成：" + newname)
	if ext == VideoTypeM4s {
		err = merge.Add(newname, video.Cid)
		if err != nil {
			return []*crawler.Request{merge.AddReq(newname, video.Cid)}, nil
		}

	}
	return nil, nil
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
