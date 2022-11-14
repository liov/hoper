package main

import (
	"bytes"
	"fmt"
	"github.com/actliboy/hoper/server/go/lib/initialize"
	"github.com/actliboy/hoper/server/go/lib/utils/dao/db/postgres"
	"github.com/actliboy/hoper/server/go/lib/utils/fs"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"tools/bilibili/config"
	"tools/bilibili/dao"
	"tools/bilibili/rpc"
)

func main() {
	defer initialize.Start(config.Conf, &dao.Dao)()
	//delete("F:\\B站\\video\\10139490\\10139490_207568591_395475557_～Alone～_alone_1_120.flv")
	deduplication()
}

func fixRecord() {
	dir := "F:\\B站\\"
	var videos []*Video
	dao.Dao.Hoper.DB.Raw(`SELECT a.owner->'mid' up_id,b.aid,b.cid,a.title,a.p->'page' page,a.p->'part' part
FROM ` + dao.TableNameVideo + ` b 
LEFT JOIN (SELECT data->'title' title, data->'owner' owner, jsonb_path_query(data,'$.pages[*]') p FROM ` + dao.TableNameView + `)  a ON (a.p->'cid')::int8 = b.cid
WHERE b.created_at > '2022-09-28 11:33:26' AND b.record = 0 AND b.` + postgres.NotDeleted + ` ORDER BY created_at`).Find(&videos)

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

	commondir := "D:\\F\\B站\\"

	videoPath := make(map[int]*VideoName)
	picPath := make(map[int]string)
	vdir := commondir + "video_bak"
	var cids []int
	files, _ := os.ReadDir(vdir)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		strs := strings.Split(file.Name(), "_")
		aid, _ := strconv.Atoi(strs[0])
		cid, _ := strconv.Atoi(strs[1])
		cids = append(cids, cid)
		videoPath[cid] = &VideoName{aid, file.Name(), strs}
		if len(cids) == 100 {
			fixNameVideoHelper(vdir, cids, videoPath)
			cids = cids[:0]
			videoPath = make(map[int]*VideoName)
		}
	}
	if len(cids) > 0 {
		fixNameVideoHelper(vdir, cids, videoPath)
		cids = cids[:0]
		videoPath = make(map[int]*VideoName)
	}
	fmt.Println(len(videoPath))
	pdir := commondir + "pic_bak"
	files, _ = os.ReadDir(pdir)
	var aids []int
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		aidStr := strings.Split(file.Name(), "_")[0]
		aid, _ := strconv.Atoi(aidStr)
		aids = append(aids, aid)
		picPath[aid] = file.Name()
		if len(aids) == 100 {
			fixNamePicHelper(pdir, aids, picPath)
			aids = aids[:0]
			picPath = make(map[int]string)
		}
	}
	fmt.Println(len(picPath))
	if len(aids) > 0 {
		fixNamePicHelper(pdir, aids, picPath)
		aids = aids[:0]
		picPath = make(map[int]string)
	}

}

func fixNameVideoHelper(vdir string, cids []int, videoPath map[int]*VideoName) {
	var videos []*Video
	dao.Dao.Hoper.DB.Raw(`SELECT a.owner->'mid' up_id,b.aid,b.cid,a.title,a.p->'page' page,a.p->'part' part
FROM `+dao.TableNameVideo+` b 
LEFT JOIN (SELECT data->'title' title, data->'owner' owner, jsonb_path_query(data,'$.pages[*]') p FROM `+dao.TableNameView+`)  a ON (a.p->'cid')::int8 = b.cid
WHERE b.cid IN (?) `, cids).Scan(&videos)
	for _, v := range videos {
		cv, ok := videoPath[v.Cid]
		if ok {
			dir := vdir + fs.PathSeparator + strconv.Itoa(v.UpId)
			_, err := os.Stat(dir)
			if os.IsNotExist(err) {
				err = os.Mkdir(dir, 0666)
				if err != nil {
					log.Println(err)
				}
			}
			part := fs.PathClean(v.Part)

			if strings.HasSuffix(cv.strs[2], part) {
				part = "!part=title!"
			}
			newpath := dir + fs.PathSeparator + fmt.Sprintf("%d_%s_%s_%s_%s_%s_%s", v.UpId, cv.strs[0], cv.strs[1], cv.strs[2], part, cv.strs[3], cv.strs[4])
			fmt.Println("rename:", newpath)
			err = os.Rename(vdir+fs.PathSeparator+cv.VideoFileName, newpath)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

type PicName struct {
	Aid  int
	UpId int
}

func fixNamePicHelper(pdir string, aids []int, picPath map[int]string) {
	var pics []*PicName
	dao.Dao.Hoper.DB.Raw(`SELECT a.owner->'mid' up_id,a.aid
FROM  (SELECT  data->'owner' owner, aid FROM `+dao.TableNameView+` WHERE aid IN (?))  a`, aids).Scan(&pics)
	for _, v := range pics {
		cp, ok := picPath[v.Aid]
		if ok {
			dir := pdir + fs.PathSeparator + strconv.Itoa(v.UpId)
			_, err := os.Stat(dir)
			if os.IsNotExist(err) {
				err = os.Mkdir(dir, 0666)
				if err != nil {
					log.Println(err)
				}
			}
			newpath := dir + fs.PathSeparator + strconv.Itoa(v.UpId) + "_" + cp
			fmt.Println("rename:", newpath)
			err = os.Rename(pdir+fs.PathSeparator+cp, newpath)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func rename() {
	commondir := "F:\\B站\\"
	dirs, _ := os.ReadDir(commondir)
	for _, dir := range dirs {
		if dir.IsDir() {
			newpicDir := commondir + "pic\\" + dir.Name() + "\\"
			_, err := os.Stat(newpicDir)
			if os.IsNotExist(err) {
				err = os.Mkdir(newpicDir, 0666)
				if err != nil {
					log.Println(err)
				}
			}
			newvidDir := commondir + "video\\" + dir.Name() + "\\"
			_, err = os.Stat(newvidDir)
			if os.IsNotExist(err) {
				err = os.Mkdir(newvidDir, 0666)
				if err != nil {
					log.Println(err)
				}
			}
			subpicdir := commondir + dir.Name() + "\\pic\\"
			files, _ := os.ReadDir(subpicdir)
			for _, file := range files {
				fileName := file.Name()
				filePath := subpicdir + file.Name()
				err = os.Rename(filePath, newpicDir+fileName)
				if err != nil {
					log.Println(err)
				}
			}
			files, _ = os.ReadDir(subpicdir)
			if len(files) == 0 {
				err = os.Remove(subpicdir)
				if err != nil {
					log.Println(err)
				}
			}
			subviddir := commondir + dir.Name() + "\\video\\"
			files, _ = os.ReadDir(subviddir)
			for _, file := range files {
				fileName := file.Name()
				filePath := subviddir + file.Name()
				err = os.Rename(filePath, newvidDir+fileName)
				if err != nil {
					log.Println(err)
				}
			}
			files, _ = os.ReadDir(subviddir)
			if len(files) == 0 {
				err = os.Remove(subviddir)
				if err != nil {
					log.Println(err)
				}
			}
		}

		files, _ := os.ReadDir(commondir + dir.Name())
		if len(files) == 0 {
			err := os.Remove(commondir + dir.Name())
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func delete(name string) {
	strs := strings.Split(name, "_")
	cid := strs[2]
	dao.Dao.Hoper.Table(dao.TableNameVideo).Where("cid = "+cid).Update("record", 0)
}

var apiservice = &rpc.API{}

func deduplication() {
	dir := "F:\\B站\\video"
	log.Println(path.Dir(dir))
	dirs, _ := os.ReadDir(dir)

	for _, subdir := range dirs {
		if subdir.IsDir() {
			m := make(map[string]string)
			files, _ := os.ReadDir(dir + fs.PathSeparator + subdir.Name())
			for _, file := range files {
				cid := strings.Split(file.Name(), "_")[2]
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
					os.Remove(dir + fs.PathSeparator + subdir.Name() + fs.PathSeparator + remove)
				} else {
					m[cid] = file.Name()
				}

			}
		}
	}
}
