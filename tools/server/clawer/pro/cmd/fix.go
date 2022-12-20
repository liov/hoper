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
	rename()
}

func rename() {
	dir := "F:\\Pictures\\pron\\91\\pic_5"
	log.Println(path.Dir(dir))
	dirs, _ := os.ReadDir(dir)
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
							allSubDir3 := allSubDir2 + fs.PathSeparator + subdir3.Name()

							files, _ := os.ReadDir(allSubDir3)
							var has bool
							for _, f := range files {
								if strings.HasSuffix(f.Name(), ".txt") {
									date := strings.Replace(f.Name()[:len(f.Name())-4], " ", "-", 1)
									log.Println(allSubDir3+fs.PathSeparator+"index.html", allSubDir3+fs.PathSeparator+date+".html")
									os.Rename(allSubDir3+fs.PathSeparator+"index.html", allSubDir3+fs.PathSeparator+date+".html")
									os.Rename(allSubDir3+fs.PathSeparator+f.Name(), allSubDir3+fs.PathSeparator+date+".txt")
									has = true
									break
								}
							}
							if !has {
								opath := "E" + allSubDir3[1:]
								files, _ := os.ReadDir(opath)
								for _, f := range files {
									if strings.Contains(f.Name(), " ") {
										date := strings.Replace(f.Name(), " ", "-", 1)
										log.Println(allSubDir3+fs.PathSeparator+"index.html", allSubDir3+fs.PathSeparator+date+".html")
										os.Rename(allSubDir3+fs.PathSeparator+"index.html", allSubDir3+fs.PathSeparator+date+".html")
										has = true
										break
									}
								}
							}
						}

					}
				}
			}
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
