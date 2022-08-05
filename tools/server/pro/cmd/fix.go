package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/actliboy/hoper/server/go/lib/utils/fs"
	"tools/pro"
)

func main() {
	//pro.SetDB()
	pro.Start(fixOne)
}

func fixOne(sd *pro.Speed) {
	fixPic(`fail_pic_2022_01_24_16_27_50`, sd)
}

func fix(sd *pro.Speed) {
	fileInfos, err := ioutil.ReadDir(pro.Conf.Pro.CommonDir)
	if err != nil {
		log.Println(err)
	}
	for i := range fileInfos {
		if !fileInfos[i].IsDir() && strings.HasPrefix(fileInfos[i].Name(), "fail_post") {
			pro.FixWeb(fileInfos[i].Name(), sd, pro.Fetch)

		}
	}
	fileInfos, err = ioutil.ReadDir(pro.Conf.Pro.CommonDir)
	if err != nil {
		log.Println(err)
	}
	for i := range fileInfos {
		if !fileInfos[i].IsDir() && strings.HasPrefix(fileInfos[i].Name(), "fail_pic") {
			fixPic(fileInfos[i].Name(), sd)
		}
	}
}

func fixPic(path string, sd *pro.Speed) {
	f, err := os.Open(pro.Conf.Pro.CommonDir + path + pro.Conf.Pro.Ext)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), "<->")
		img, dir := s[0], s[1]
		dir = fs.PathEdit(dir)
		log.Println(img, dir)
		sd.Add()
		go pro.Download(img, dir, sd)
		time.Sleep(pro.Conf.Pro.Interval)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
