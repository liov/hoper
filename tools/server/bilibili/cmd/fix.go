package main

import (
	"context"
	"fmt"
	"github.com/actliboy/hoper/server/go/lib/initialize"
	"github.com/actliboy/hoper/server/go/lib/utils/dao/db/postgres"
	"github.com/actliboy/hoper/server/go/lib/utils/fs"
	osi "github.com/actliboy/hoper/server/go/lib/utils/os"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"tools/bilibili/config"
	"tools/bilibili/dao"
	"tools/bilibili/download"
	"tools/bilibili/rpc"
	"tools/bilibili/tool"
)

func main() {
	defer initialize.Start(config.Conf, &dao.Dao)()
	rename()
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

func fixName() {
	type VideoName struct {
		Aid           int
		VideoFileName string
		strs          []string
	}

	commondir := "G:\\B站\\"

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
LEFT JOIN (SELECT data->'title' title, data->'owner' owner, jsonb_path_query(data,'$.pages[*]') p FROM ` + dao.TableNameView + `)  a ON (a.p->'cid')::int8 = b.cid
WHERE b.` + postgres.NotDeleted + ` LIMIT 100 OFFSET ` + strconv.Itoa(pageNo*pageSize)).Find(&videos)

		for _, v := range videos {
			cv, ok := videoPath[v.Cid]
			if ok {
				dir := commondir
				_, err := os.Stat(dir)
				if os.IsNotExist(err) {
					err = os.Mkdir(dir, 0666)
					if err != nil {
						log.Println(err)
					}
				}
				_, err = os.Stat(dir + "video\\" + strconv.Itoa(v.UpId))
				if os.IsNotExist(err) {
					err = os.Mkdir(dir+"video\\"+strconv.Itoa(v.UpId), 0666)
					if err != nil {
						log.Println(err)
					}
				}
				part := fs.PathClean(v.Part)
				if part == cv.strs[2] {
					part = "!part=title!"
				}
				newpath := dir + "\\video\\" + strconv.Itoa(v.UpId) + "\\" + fmt.Sprintf("%d_%s_%s_%s_%s_%s_%s", v.UpId, cv.strs[0], cv.strs[1], cv.strs[2], part, cv.strs[3], cv.strs[4])
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
				_, err = os.Stat(dir + "\\pic\\" + strconv.Itoa(v.UpId))
				if os.IsNotExist(err) {
					err = os.Mkdir(dir+"\\pic\\"+strconv.Itoa(v.UpId), 0666)
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

func fixCover() {
	apiservice := rpc.API{}
	res, err := apiservice.GetFavLResourceList(63181530, 5)
	if err != nil {
		log.Println(err)
	}
	for _, fav := range res.Medias {
		aid := tool.Bv2av(fav.Bvid)
		err = download.CoverDownload(context.Background(), fav.Cover, fav.Upper.Mid, aid)
		if err != nil {
			log.Println("下载图片失败：", err)
		}
	}

}

func transferCodec() {
	command := `ffmpeg  -hwaccel_output_format qsv -c:v h264_qsv -i %s -c:v hevc_qsv -global_quality 23  -gpu_copy 1 -c:a copy %s.mp4`
	commondir := "D:\\F\\B站\\"
	dirs, _ := os.ReadDir(commondir)
	for _, dir := range dirs {
		if dir.IsDir() {
			subdir := commondir + dir.Name() + "\\video\\"
			files, _ := os.ReadDir(subdir)
			for _, file := range files {
				fileName := file.Name()
				if !strings.HasSuffix(fileName, ".flv") {
					continue
				}
				filePath := subdir + file.Name()
				command = fmt.Sprintf(command, filePath, filePath[:len(filePath)-len(".flv")])
				log.Println(command)
				res, err := osi.CMD(command)
				if err != nil {
					log.Println(err)
					continue
				}
				log.Println(res)
				//log.Println(os.Remove(filePath))
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
