package download

import (
	"context"
	"fmt"
	"github.com/actliboy/hoper/server/go/lib/utils/conctrl"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client/crawler"
	osi "github.com/actliboy/hoper/server/go/lib/utils/os"
	"log"
	"os"
	"sync"
	"tools/bilibili/config"
	"tools/bilibili/dao"
)

var merge VideoMerge

type VideoMerge struct {
	Map sync.Map
}

func (m *VideoMerge) AddReq(path string, cid int) *crawler.Request {
	return &crawler.Request{
		TaskMeta: conctrl.TaskMeta{},
		Key:      "",
		TaskFunc: func(ctx context.Context) ([]*crawler.Request, error) {
			return nil, m.Add(path, cid)
		},
	}
}

func (m *VideoMerge) Add(path string, cid int) error {
	path = path[:len(path)-len(".m4s.video")]
	if _, ok := m.Map.Load(path); ok {
		err := mergeVideo(path, cid)
		if err != nil {
			return err
		}
		m.Map.Delete(path)
	} else {
		m.Map.Store(path, true)
	}
	return nil
}

func mergeVideo(path string, cid int) error {
	cmd := config.Conf.Bilibili.FFmpegPath + fmt.Sprintf(" -i %s.m4s.video -i %s.m4s.audio -c copy -strict experimental %s.mp4", path, path, path)
	_, err := osi.CMD(cmd)
	if err != nil {
		log.Println("合并失败：", err)
		return err
	}
	os.Remove(path + ".m4s.video")
	os.Remove(path + ".m4s.audio")
	dao.Dao.Hoper.Table(dao.TableNameVideo).Where("cid = ?", cid).Update("record", 2)
	log.Println("合并完成：" + path)
	return nil
}
