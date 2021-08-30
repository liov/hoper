package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/liov/hoper/server/go/lib/utils/fs"
	"tools/pro"
)

func main() {
	//pro.SetDB()
	pro.Start(fixOne)
}

func fixOne(sd *pro.Speed) {
	fixPic(`fail_pic_2021_02_08_18_04_49`, sd)
}

func fix(sd *pro.Speed) {
	fileInfos, err := ioutil.ReadDir(pro.CommonDir)
	if err != nil {
		log.Println(err)
	}
	for i := range fileInfos {
		if !fileInfos[i].IsDir() && !strings.HasSuffix(fileInfos[i].Name(), "fail_pic") {
			pro.FixWeb(fileInfos[i].Name(), sd, pro.Fetch)

		}
	}
	fileInfos, err = ioutil.ReadDir(pro.CommonDir)
	if err != nil {
		log.Println(err)
	}
	for i := range fileInfos {
		if !fileInfos[i].IsDir() && strings.HasSuffix(fileInfos[i].Name(), "fail_pic") {
			fixPic(fileInfos[i].Name(), sd)
		}
	}
}

func fixPic(path string, sd *pro.Speed) {
	f, err := os.Open(pro.CommonDir + path + pro.Ext)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), "<->")
		img, dir := s[0], s[1]
		dir = fs.PathClean(dir)
		log.Println(img, dir)
		sd.Add(1)
		go pro.Download(img, dir, sd)
		time.Sleep(pro.Interval)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
