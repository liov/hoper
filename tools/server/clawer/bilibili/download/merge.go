package download

import (
	"context"
	"fmt"
	"github.com/liov/hoper/server/go/lib/utils/fs"
	timei "github.com/liov/hoper/server/go/lib/utils/time"
	"github.com/liov/hoper/server/go/lib_v2/utils/conctrl"
	"github.com/liov/hoper/server/go/lib_v2/utils/net/http/client/crawler"
	claweri "tools/clawer"

	osi "github.com/liov/hoper/server/go/lib/utils/os"
	"log"
	"os"
	"strconv"
	"sync"
	"tools/clawer/bilibili/config"
	"tools/clawer/bilibili/dao"
)

var merge VideoMerge

func GetEngineMerge(engine *crawler.Engine) *VideoMerge {
	merge.engine = engine
	merge.fixedWorkerId = engine.NewFixedWorker(0)
	return &merge
}

type VideoMerge struct {
	Map           sync.Map
	engine        *crawler.Engine
	fixedWorkerId int
}

func (m *VideoMerge) Add(video *Video) {
	if single, ok := m.Map.Load(video.Cid); ok {
		m.engine.BaseEngine.AddFixedTask(merge.fixedWorkerId, &conctrl.BaseTask[string, crawler.Prop]{
			BaseTaskFunc: func(ctx context.Context) {
				MergeVideo(video, single.(bool))
				m.Map.Delete(video.Cid)
			},
		})
	} else {
		m.Map.Store(video.Cid, false)
	}
}

func MergeVideo(video *Video, single bool) error {
	pubAt := video.PubAt.Format(timei.TimeFormatCompact)
	src := fmt.Sprintf("%d_%d_%d", video.Uid, video.Aid, video.Cid)
	dst := pubAt + "_" + src + "_" + video.Title + "_" + video.Part + "_" + strconv.Itoa(video.Quality)
	fpath := config.Conf.Bilibili.DownloadTmpPath + fs.PathSeparator + src
	dir := config.Conf.Bilibili.DownloadPath + fs.PathSeparator + strconv.Itoa(video.Uid) + fs.PathSeparator + pubAt[:4]

	os.MkdirAll(dir, 0666)
	var ext string
	if video.CodecId == VideoTypeM4sCodec12 {
		ext = ".mp4"
	} else {
		ext = ".flv"
	}
	mergePath := dir + fs.PathSeparator + src + ext
	renamePath := dir + fs.PathSeparator + dst + ext
	// 开发过程的bug，这里兼容解决一下，都先检查模板文件是否存在，不存在才执行响应操作
	if fs.NotExist(renamePath) {
		if fs.NotExist(mergePath) {
			var cmd string
			if single {
				cmd = config.Conf.Bilibili.FFmpegPath + fmt.Sprintf(" -i %s.m4s.video  -c copy -strict experimental %s", fpath, mergePath)
			} else {
				cmd = config.Conf.Bilibili.FFmpegPath + fmt.Sprintf(" -i %s.m4s.video -i %s.m4s.audio -c copy -strict experimental %s", fpath, fpath, mergePath)
			}
			_, err := osi.CMD(cmd)
			if err != nil {
				log.Println("合并失败：", dst, err)
				log.Println("cmd:", cmd)
				return err
			}
		}

		err := os.Rename(mergePath, dir+fs.PathSeparator+dst+ext)
		if err != nil {
			log.Println(err)
			return err
		}

	}

	err := os.Remove(fpath + ".m4s.video")
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
	cdir := claweri.Dir{
		Platform: 3,
		UserId:   video.Uid,
		KeyId:    video.Cid,
		KeyIdStr: fmt.Sprintf("%d_%d", video.Aid, video.Cid),
		BaseUrl:  fmt.Sprintf("%s_%s_%d", video.Title, video.Part, video.Quality) + ext,
		Type:     3,
		PubAt:    video.PubAt,
	}
	dao.Dao.Hoper.Create(&cdir)
	log.Println("合并完成：" + dst)
	return nil
}
