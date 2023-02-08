package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/liov/hoper/server/go/lib/initialize"
	dbi "github.com/liov/hoper/server/go/lib/utils/dao/db/const"
	"github.com/liov/hoper/server/go/lib/utils/fs"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
	claweri "tools/clawer"
	"tools/clawer/bilibili/config"
	"tools/clawer/bilibili/dao"
	"tools/clawer/bilibili/download"
	"tools/clawer/bilibili/rpc"
)

func main() {
	defer initialize.Start(config.Conf, &dao.Dao)()
	//delete("F:\\B站\\video\\10139490\\10139490_207568591_395475557_～Alone～_alone_1_120.flv")
	deduplication()
	//fixName()
	/*fs.RangeDir("F:\\debug\\B站", func(dir string, entry os.DirEntry) error {
		if strings.HasSuffix(entry.Name(), "64.flv") {
			log.Println(entry.Name())
			os.Remove(dir + fs.PathSeparator + entry.Name())
		}
		return nil
	})*/
}

func fixRecord() {
	dir := "F:\\B站\\"
	var videos []*Video
	dao.Dao.Hoper.DB.Raw(`SELECT a.owner->'mid' up_id,b.aid,b.cid,a.title,a.p->'page' page,a.p->'part' part
FROM ` + dao.TableNameVideo + ` b 
LEFT JOIN (SELECT data->'title' title, data->'owner' owner, jsonb_path_query(data,'$.pages[*]') p FROM ` + dao.TableNameView + `)  a ON (a.p->'cid')::int8 = b.cid
WHERE b.created_at > '2022-09-28 11:33:26' AND b.record = 0 AND b.` + dbi.NotDeleted + ` ORDER BY created_at`).Find(&videos)

	for _, video := range videos {
		if video.Part != video.Title {
			continue
		}
		for _, q := range []int{112, 116, 120} {
			filename := fmt.Sprintf("%d_%d_%d_%s_%s_%d_%d.flv", video.UpId, video.Aid, video.Cid, video.Title, video.Part, 1, q)
			oldpath := dir + strconv.Itoa(video.UpId) + "\\video\\" + fs.PathClean(filename)
			info, err := os.Stat(oldpath)
			if err != nil {
				log.Println(err)
				continue
			}
			if info != nil {
				dao.Dao.Hoper.Table(dao.TableNameVideo).Where("cid = ?", video.Cid).Update("record", 1)
			}
		}
	}
}

func fixRecord2() {
	dir := "D:\\F\\video\\"

	dirs, _ := os.ReadDir(dir)
	for _, subdir := range dirs {
		files, _ := os.ReadDir(dir + subdir.Name())
		for _, file := range files {
			if strings.HasSuffix(file.Name(), "64.flv") {
				cid := strings.Split(file.Name(), "_")[2]
				dao.Dao.Hoper.Table(dao.TableNameVideo).Where("cid = "+cid).Update("record", 1)
			}
		}
	}
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
				dao.Dao.Hoper.Table(dao.TableNameVideo).Where(`cid = `+cid).UpdateColumn("record", 1)
				os.Remove(path.Join(dir, file.Name()))
			}
		}
	}
}

func remove() {
	dir := "F:\\B站\\video"
	log.Println(path.Dir(dir))
	dirs, _ := os.ReadDir(dir)
	m := map[string]struct{}{}
	for _, subdir := range dirs {
		if subdir.IsDir() {
			files, _ := os.ReadDir(dir + fs.PathSeparator + subdir.Name())
			for _, file := range files {
				cid := strings.Split(file.Name(), "_")[2]
				m[cid] = struct{}{}
				/*err := dao.Dao.Hoper.Table(dao.TableNameVideo).Where("cid = "+cid).Update("record", true).Error
				if err != nil {
					log.Println(err)
					return
				}*/
			}
		}
	}
	dir = "F:\\B站"
	files, _ := os.ReadDir(dir)
	var buf bytes.Buffer
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		var cid string
		if strings.Contains(file.Name(), "-") {
			cid = strings.Split(file.Name(), "-")[0]
		}
		if strings.Contains(file.Name(), "_") {
			cid = strings.Split(file.Name(), "_")[0]
		}
		buf.WriteString(cid)
		buf.WriteString(",")
		//log.Println(cid)
		if _, ok := m[cid]; ok {
			err := os.Remove(path.Join(dir, file.Name()))
			if err != nil {
				log.Println(err)
				return
			}
			log.Println("remove", file.Name())
		}
	}
	log.Println(buf.String())
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

type VideoName struct {
	Aid           int
	VideoFileName string
	strs          []string
}

func fixName() {

	commondir := "F:\\B站\\0\\0001"
	timer := time.NewTicker(time.Second)
	files, _ := os.ReadDir(commondir)
	uid := 0
	m := make(map[int]time.Time)
	ctx := context.Background()
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		info, _ := file.Info()
		strs := strings.Split(file.Name(), "_")
		aid, _ := strconv.Atoi(strs[0])
		pubAt, ok := m[aid]
		if !ok {
			var view dao.View
			err := dao.Dao.Hoper.Where(`aid = ` + strs[0]).First(&view).Error
			if err != nil {
				log.Println(err)
				<-timer.C
				view2, err := download.RecordViewInfo(ctx, aid)
				if view2 == nil {
					log.Println(err)
					continue
				}
				for _, page := range view2.Pages {
					if len(view2.Pages) == 1 {
						page.Part = download.PartEqTitle
					}
					video := download.NewVideo(view2.Owner.Mid, view2.Title, view2.Aid, page.Cid, page.Page, page.Part, view2.PubDate)
					_, err := video.RecordVideo(ctx)
					if err != nil {
						log.Println(err)
						continue
					}
				}
				view = dao.View{
					Bvid:    view2.Bvid,
					Aid:     view2.Aid,
					Uid:     view2.Owner.Mid,
					Title:   view2.Title,
					Desc:    view2.Desc,
					Dynamic: view2.Dynamic,
					Tid:     view2.Tid,
					Pic:     view2.Pic,
					Ctime:   time.Unix(int64(view2.Ctime), 0),
					Tname:   view2.Tname,
					Videos:  view2.Videos,
					Pubdate: time.Unix(int64(view2.PubDate), 0),
				}
			}
			pubAt = view.Pubdate
			m[aid] = pubAt
			uid = view.Uid
		}

		dir := claweri.Dir{
			Platform:  3,
			UserId:    uid,
			KeyId:     aid,
			BaseUrl:   file.Name()[len(strs[0])+1:],
			Type:      1,
			PubAt:     pubAt,
			CreatedAt: info.ModTime(),
		}

		newpath := "F:\\B站\\" + dir.Path()
		log.Println("rename:", commondir+fs.PathSeparator+file.Name(), newpath)
		os.MkdirAll(fs.GetDir(newpath), 0666)
		os.Rename(commondir+fs.PathSeparator+file.Name(), newpath)
		//dao.Dao.Hoper.Create(&dir)
	}
	files, _ = os.ReadDir(commondir)
	if len(files) == 0 {
		err := os.Remove(commondir)
		if err != nil {
			log.Println(err)
		}
	}

}

type PicName struct {
	Aid  int
	UpId int
}

func delete(name string) {
	strs := strings.Split(name, "_")
	cid := strs[2]
	dao.Dao.Hoper.Table(dao.TableNameVideo).Where("cid = "+cid).Update("record", 0)
}

var apiservice = &rpc.API{}

func deduplication() {
	dir := "G:\\B站"
	log.Println(path.Dir(dir))
	userDirs, _ := os.ReadDir(dir)

	for _, userDir := range userDirs {
		if userDir.IsDir() {
			subDir := dir + fs.PathSeparator + userDir.Name()
			yearDirs, _ := os.ReadDir(subDir)
			for _, yearDir := range yearDirs {
				subDir1 := subDir + fs.PathSeparator + yearDir.Name()
				m := make(map[string]string)
				files, _ := os.ReadDir(subDir1)
				for _, file := range files {
					if !strings.HasSuffix(file.Name(), ".flv") && !strings.HasSuffix(file.Name(), ".mp4") {
						continue
					}
					if strings.HasSuffix(file.Name(), "downloading") {
						os.Remove(subDir1 + fs.PathSeparator + file.Name())
						continue
					}
					cid := strings.Split(file.Name(), "_")[3]
					if path, ok := m[cid]; ok {
						remove := file.Name()
						if strings.HasSuffix(path, ".flv") {
							if strings.HasSuffix(remove, ".flv") {
								if strings.Contains(path, "_1_") {
									remove = path
								}
							} else {
								remove = path
							}

						}
						log.Println(remove)
						os.Remove(subDir1 + fs.PathSeparator + remove)
					} else {
						m[cid] = file.Name()
					}

				}
			}
		}
	}
}
