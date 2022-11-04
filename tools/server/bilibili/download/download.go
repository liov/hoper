package download

import (
	"context"
	"fmt"
	"github.com/actliboy/hoper/server/go/lib/utils/fs"
	"github.com/actliboy/hoper/server/go/lib/utils/generics/net/http/client/crawler"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client"
	"strings"

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
	UpId      int
	Title     string
	Aid       int
	Cid       int
	Page      int
	Part      string
	Quality   int
	CreatedAt time.Time
	Record    int
	CodecId   int
}

func NewVideo(upId int, title string, aid, cid, page int, part string, created time.Time) *Video {
	return &Video{
		UpId:      upId,
		Title:     title,
		Aid:       aid,
		Cid:       cid,
		Page:      page,
		Part:      part,
		CreatedAt: created,
	}
}

func (video *Video) DownloadVideoReq(typ string, order int, url string) *crawler.Request {
	return &crawler.Request{
		TaskMeta: crawler.TaskMeta{BaseTaskMeta: crawler.BaseTaskMeta{Key: "下载视频：" + strconv.Itoa(video.Cid) + typ}, Kind: KindDownloadVideo},
		TaskFunc: func(ctx context.Context) ([]*crawler.Request, error) {
			return video.DownloadVideo(typ, order, url)
		},
	}
}

func (video *Video) DownloadVideo(typ string, order int, url string) ([]*crawler.Request, error) {

	var filename string

	video.Part = fs.PathClean(video.Part)
	video.Title = fs.PathClean(video.Title)
	if strings.HasSuffix(video.Title, video.Part) {
		video.Part = PartEqTitle
	}

	if video.CodecId == VideoTypeFlv {
		filename = fmt.Sprintf("%d_%d_%d_%s_%s_%d_%d.flv.downloading", video.UpId, video.Aid, video.Cid, video.Title, video.Part, order, video.Quality)
		filename = filepath.Join(config.Conf.Bilibili.DownloadVideoPath, strconv.Itoa(video.UpId), filename)
	} else {
		filename = fmt.Sprintf("%d_%d_%d.m4s.%s.downloading", video.UpId, video.Aid, video.Cid, typ)
		filename = config.Conf.Bilibili.DownloadTmpPath + fs.PathSeparator + filename
	}

	newname := filename[:len(filename)-len(DownloadingExt)]

	_, err := os.Stat(newname)
	if os.IsNotExist(err) {
		referer := rpc.GetViewUrl(video.Aid)
		referer = referer + fmt.Sprintf("/?p=%d", video.Page)

		c := http.Client{CheckRedirect: genCheckRedirectfun(referer)}

		request, err := http.NewRequest(http.MethodGet, url, nil)
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

		file, err := fs.Create(filename)
		if err != nil {
			log.Println("错误信息：", err)
			return nil, err
		}

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
	}
	log.Println("下载完成：" + newname)

	if video.CodecId == VideoTypeFlv {
		dao.Dao.Hoper.Table(dao.TableNameVideo).Where("cid = ?", video.Cid).Update("record", 3)
	}

	if video.CodecId == VideoTypeM4sCodec12 || video.CodecId == VideoTypeM4sCodec7 {
		merge.Add(video)
	}
	return nil, nil
}

func genCheckRedirectfun(referer string) func(req *http.Request, via []*http.Request) error {
	return func(req *http.Request, via []*http.Request) error {
		req.Header.Set("Referer", referer)
		return nil
	}
}
