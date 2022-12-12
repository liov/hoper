package main

import (
	"bufio"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/liov/hoper/server/go/lib/utils/fs"
	"tools/clawer/pro"
)

func main() {
	//pro.SetDB()
	maxId()
}

func fixOne() {
	fixPic(`fail_pic_2022_01_24_16_27_50`)
}

func fix() {
	fileInfos, err := os.ReadDir(pro.Conf.Pro.CommonDir)
	if err != nil {
		log.Println(err)
	}
	for i := range fileInfos {
		if !fileInfos[i].IsDir() && strings.HasPrefix(fileInfos[i].Name(), "fail_post") {
			pro.FixWeb(fileInfos[i].Name(), pro.GetFetchReq)

		}
	}
	fileInfos, err = os.ReadDir(pro.Conf.Pro.CommonDir)
	if err != nil {
		log.Println(err)
	}
	for i := range fileInfos {
		if !fileInfos[i].IsDir() && strings.HasPrefix(fileInfos[i].Name(), "fail_pic") {
			fixPic(fileInfos[i].Name())
		}
	}
}

func fixPic(path string) {
	f, err := os.Open(pro.Conf.Pro.CommonDir + path + pro.Conf.Pro.Ext)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), "<->")
		img, dir := s[0], s[1]
		dir = fs.FileNameEdit(dir)
		log.Println(img, dir)

		go pro.Download(img, dir)
		time.Sleep(pro.Conf.Pro.Interval)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func maxId() {
	dir := "F:\\Pictures\\pron\\91\\pic_5"
	log.Println(path.Dir(dir))
	dirs, _ := os.ReadDir(dir)
	var maxId int
	for _, subdir := range dirs {
		if subdir.IsDir() {
			allSubDir := dir + fs.PathSeparator + subdir.Name()
			subDirs, _ := os.ReadDir(allSubDir)
			for _, subdir2 := range subDirs {
				if subdir2.IsDir() {
					allSubDir2 := allSubDir + fs.PathSeparator + subdir2.Name()
					subDirs2, _ := os.ReadDir(allSubDir2)
					for _, subdir3 := range subDirs2 {
						if subdir3.IsDir() {
							strs := strings.Split(subdir3.Name(), "_")
							id, _ := strconv.Atoi(strs[len(strs)-1])
							if id > maxId {
								maxId = id
							}
						}
					}
				}
			}
		}
	}
	log.Println(maxId)
}
