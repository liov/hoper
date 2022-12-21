package main

import (
	"github.com/liov/hoper/server/go/lib/initialize"
	"github.com/liov/hoper/server/go/lib/utils/fs"
	"log"
	"os"
	"strconv"
	"strings"
	"tools/clawer/timepill"
)

func main() {
	defer initialize.Start(&timepill.Conf, &timepill.Dao)()
	rename()
}

func rename() {
	commondir := "D:\\F\\timepill\\debug\\2022\\2022-12\\2022-12-21"
	files, _ := os.ReadDir(commondir)
	for _, f := range files {
		strs := strings.Split(f.Name(), "_")
		userId, _ := strconv.Atoi(strs[0])
		num := userId / 10000
		newdir := "D:\\F\\timepill\\debug\\" + strconv.Itoa(num) + "-" + strconv.Itoa(num+1) + fs.PathSeparator + strs[0] + fs.PathSeparator + "2022"
		os.MkdirAll(newdir, 0666)
		log.Println("rename:", commondir+fs.PathSeparator+f.Name(), newdir+fs.PathSeparator+strings.Join([]string{strs[2], strs[0], strs[1], strs[3]}, "_"))
		os.Rename(commondir+fs.PathSeparator+f.Name(), newdir+fs.PathSeparator+strings.Join([]string{strs[2], strs[0], strs[1], strs[3]}, "_"))
	}

}

/*func rename2() {
	commondir := "D:\\F\\timepill\\2010_"
	subdirs, _ := os.ReadDir(commondir)

	for _, subdir := range subdirs {
		compsubdir := commondir + fs.PathSeparator + subdir.Name()
		m := make(map[string]time.Time)
		files, _ := os.ReadDir(compsubdir)
		for _, f := range files {
			fname := f.Name()
			info, _ := f.Info()
			strs := strings.Split(fname, "_")

			for i := 0; i < 10; i++ {

			}
		}

		userId, _ := strconv.Atoi(strs[0])
		dir := &claweri.Dir{
			Platform:  4,
			UserId:    userId,
			KeyIdStr:  strs[1],
			BaseUrl:   strs[2],
			Type:      1,
			PubAt:     date,
			CreatedAt: info.ModTime(),
		}
		if strings.HasSuffix(strs[2], ".mov") {
			dir.Type = 2
		}
		oldpath := compsubdir + fs.PathSeparator + fname
		newpath := timepill.Conf.TimePill.PhotoPath + "/" + dir.Path()
		os.MkdirAll(fs.GetDir(newpath), 0666)
		log.Println("rename:", oldpath, newpath)
		err := os.Rename(oldpath, newpath)
		if err != nil {
			log.Println(err)
		}
		timepill.Dao.Hoper.Create(dir)
	}
}
*/
