package main

import (
	"fmt"
	"github.com/actliboy/hoper/server/go/lib/initialize"
	"github.com/actliboy/hoper/server/go/lib/utils/dao/db/postgres"
	"github.com/actliboy/hoper/server/go/lib/utils/fs"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client/crawler"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"tools/bilibili/config"
	"tools/bilibili/dao"
	"tools/bilibili/download"
)

func main() {
	defer initialize.Start(config.Conf, &dao.Dao)()
	fixName()
}

func fixRecord() {
	dir := "F:\\B站\\pic"
	files, _ := os.ReadDir(dir)
	var reqs []*crawler.Request
	for _, file := range files {
		aidStr := strings.Split(file.Name(), "_")[0]
		aid, _ := strconv.Atoi(aidStr)
		req := download.GetViewInfoReqV2(aid)
		reqs = append(reqs, req)
	}
	engine := crawler.New(10).SkipKind(4)
	engine.Run(reqs...)
}

func fixQuality() {
	dir := "D:\\F\\B站\\video"
	files, _ := os.ReadDir(dir)
	for _, file := range files {
		if strings.HasSuffix(file.Name(), "64.flv") || strings.HasSuffix(file.Name(), "80.flv") {
			cid := strings.Split(file.Name(), "_")[1]
			var quality string
			err := dao.Dao.Hoper.Raw(`SELECT data #> '{accept_quality,0}' quality FROM ` + dao.TableNameVideo + ` WHERE cid = ` + cid).Row().Scan(&quality)
			if err != nil {
				log.Println(err)
				return
			}
			fmt.Println(file.Name()[len(file.Name())-6 : len(file.Name())-4])
			if quality != file.Name()[len(file.Name())-6:len(file.Name())-4] {
				dao.Dao.Hoper.Table(dao.TableNameVideo).Where(`cid = `+cid).UpdateColumn("record", false)
				os.Remove(path.Join(dir, file.Name()))
			}
		}
	}
}

func remove() {
	dir := "F:\\B站\\video"
	log.Println(path.Dir(dir))
	files, _ := os.ReadDir(dir)
	m := map[string]struct{}{}
	for _, file := range files {
		cid := strings.Split(file.Name(), "_")[1]
		m[cid] = struct{}{}
		/*err := dao.Dao.Hoper.Table(dao.TableNameVideo).Where("cid = "+cid).Update("record", true).Error
		if err != nil {
			log.Println(err)
			return
		}*/
	}
	dir = "F:\\Pictures\\B站"
	files, _ = os.ReadDir(dir)
	for _, file := range files {
		if strings.Contains(file.Name(), "-") {
			cid := strings.Split(file.Name(), "-")[0]
			if _, ok := m[cid]; ok {
				err := os.Remove(path.Join(dir, file.Name()))
				if err != nil {
					log.Println(err)
					return
				}
				log.Println("remove", file.Name())
			}
		}
	}
}

func fixName() {
	type VideoName struct {
		Aid           int
		VideoFileName string
		strs          []string
	}
	type Video struct {
		UpId    int
		Title   string
		Aid     int
		Cid     int
		Page    int
		Part    string
		Order   int
		Quality int
	}

	commondir := "F:\\B站\\"

	videoPath := make(map[int]*VideoName)
	picPath := make(map[int]string)
	vdir := commondir + "video"
	files, _ := os.ReadDir(vdir)
	for _, file := range files {
		strs := strings.Split(file.Name(), "_")
		aid, _ := strconv.Atoi(strs[0])
		cid, _ := strconv.Atoi(strs[1])
		videoPath[cid] = &VideoName{aid, file.Name(), strs}
	}
	fmt.Println(len(videoPath))
	pdir := commondir + "pic"
	files, _ = os.ReadDir(pdir)
	for _, file := range files {
		aidStr := strings.Split(file.Name(), "_")[0]
		aid, _ := strconv.Atoi(aidStr)
		picPath[aid] = file.Name()
	}
	fmt.Println(len(picPath))
	pageNo, pageSize := 0, 100
	for {
		var videos []*Video
		dao.Dao.Hoper.DB.Raw(`SELECT a.owner->'mid' up_id,b.aid,b.cid,a.title,a.p->'page' page,a.p->'part' part
FROM ` + dao.TableNameVideo + ` b 
LEFT JOIN (SELECT data->'title' title , data->'owner' owner ,jsonb_path_query(data,'$.pages[*]') p FROM ` + dao.TableNameView + `)  a ON (a.p->'cid')::int8 = b.cid
WHERE b.` + postgres.NotDeleted + ` LIMIT 100 OFFSET ` + strconv.Itoa(pageNo*pageSize)).Find(&videos)

		for _, v := range videos {
			cv, ok := videoPath[v.Cid]
			if ok {
				dir := commondir + strconv.Itoa(v.UpId)
				_, err := os.Stat(dir)
				if os.IsNotExist(err) {
					err = os.Mkdir(dir, 0666)
					if err != nil {
						log.Println(err)
					}
				}
				_, err = os.Stat(dir + "\\" + config.Conf.Bilibili.DownloadVideoPath)
				if os.IsNotExist(err) {
					err = os.Mkdir(dir+"\\"+config.Conf.Bilibili.DownloadVideoPath, 0666)
					if err != nil {
						log.Println(err)
					}
				}
				part := fs.PathClean(v.Part)
				if part == cv.strs[2] {
					part = "!part=title!"
				}
				newpath := dir + "\\video\\" + fmt.Sprintf("%d_%s_%s_%s_%s_%s_%s", v.UpId, cv.strs[0], cv.strs[1], cv.strs[2], part, cv.strs[3], cv.strs[4])
				fmt.Println("rename:", newpath)
				err = os.Rename(vdir+"\\"+cv.VideoFileName, newpath)
				if err != nil {
					log.Println(err)
				}
			}
			cp, ok := picPath[v.Aid]
			if ok {
				dir := commondir + strconv.Itoa(v.UpId)
				_, err := os.Stat(dir)
				if os.IsNotExist(err) {
					err = os.Mkdir(dir, 0666)
					if err != nil {
						log.Println(err)
					}
				}
				_, err = os.Stat(dir + "\\" + config.Conf.Bilibili.DownloadPicPath)
				if os.IsNotExist(err) {
					err = os.Mkdir(dir+"\\"+config.Conf.Bilibili.DownloadPicPath, 0666)
					if err != nil {
						log.Println(err)
					}
				}
				_, err = os.Stat(pdir + "\\" + cp)
				if os.IsNotExist(err) {
					continue
				}
				newpath := dir + "\\pic\\" + strconv.Itoa(v.UpId) + "_" + cp
				fmt.Println("rename:", newpath)
				err = os.Rename(pdir+"\\"+cp, newpath)
				if err != nil {
					log.Println(err)
				}
			}
		}
		if len(videos) < pageSize {
			break
		}
		pageNo++
	}
}
