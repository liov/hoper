package download

import (
	"context"
	"fmt"
	"github.com/actliboy/hoper/server/go/lib/utils/conctrl"
	"github.com/actliboy/hoper/server/go/lib/utils/fs"
	"github.com/actliboy/hoper/server/go/lib/utils/generics/net/http/client/crawler"

	osi "github.com/actliboy/hoper/server/go/lib/utils/os"
	"log"
	"os"
	"strconv"
	"sync"
	"tools/bilibili/config"
	"tools/bilibili/dao"
)

var merge = VideoMerge{
	ctrl: make(conctrl.Controller),
}

func GetMerge() *VideoMerge {
	return &merge
}

type VideoMerge struct {
	Map  sync.Map
	ctrl conctrl.Controller
}

func (m *VideoMerge) AddReq(video *Video) *crawler.Request {
	return &crawler.Request{
		TaskFunc: func(ctx context.Context) ([]*crawler.Request, error) {
			return nil, m.Add(video)
		},
	}
}

func (m *VideoMerge) Add(video *Video) error {
	if single, ok := m.Map.Load(video.Cid); ok {
		m.ctrl.AddTask(func() error {
			err := MergeVideo(video, single.(bool))
			if err != nil {
				return err
			}
			m.Map.Delete(video.Cid)
			return nil
		})
	} else {
		m.Map.Store(video.Cid, false)
	}
	return nil
}

func MergeVideo(video *Video, single bool) error {
	src := fmt.Sprintf("%d_%d_%d", video.UpId, video.Aid, video.Cid)
	dst := src + "_" + video.Title + "_" + video.Part + "_" + strconv.Itoa(video.Quality)

	fpath := config.Conf.Bilibili.DownloadTmpPath + fs.PathSeparator + src
	dir := config.Conf.Bilibili.DownloadVideoPath + fs.PathSeparator + strconv.Itoa(video.UpId)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		os.Mkdir(dir, 0666)
	}
	var ext string
	if video.CodecId == VideoTypeM4sCodec12 {
		ext = ".mp4"
	} else {
		ext = ".flv"
	}
	mergePath := dir + fs.PathSeparator + src + ext
	renamePath := dir + fs.PathSeparator + dst + ext
	// 开发过程的bug，这里兼容解决一下，都先检查模板文件是否存在，不存在才执行响应操作
	_, err = os.Stat(renamePath)
	if os.IsNotExist(err) {
		_, err = os.Stat(mergePath)
		if os.IsNotExist(err) {
			var cmd string
			if single {
				cmd = config.Conf.Bilibili.FFmpegPath + fmt.Sprintf(" -i %s.m4s.video  -c copy -strict experimental %s", fpath, mergePath)
			} else {
				cmd = config.Conf.Bilibili.FFmpegPath + fmt.Sprintf(" -i %s.m4s.video -i %s.m4s.audio -c copy -strict experimental %s", fpath, fpath, mergePath)
			}
			_, err = osi.CMD(cmd)
			if err != nil {
				log.Println("合并失败：", dst, err)
				log.Println("cmd:", cmd)
				return err
			}
		}

		err = os.Rename(mergePath, dir+fs.PathSeparator+dst+ext)
		if err != nil {
			log.Println(err)
			return err
		}

	}

	err = os.Remove(fpath + ".m4s.video")
	if err != nil {
		log.Println(err)
	}
	if !single {
		err = os.Remove(fpath + ".m4s.audio")
		if err != nil {
			log.Println(err)
		}
	}
	record := 2
	if video.CodecId == VideoTypeM4sCodec7 {
		record = 3
	}

	dao.Dao.Hoper.Table(dao.TableNameVideo).Where("cid = ?", video.Cid).Update("record", record)
	log.Println("合并完成：" + dst)
	return nil
}

func (m *VideoMerge) Start() {
	go m.ctrl.Start()
}
