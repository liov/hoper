package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"tools/pro"
)

func main() {
	pro.Start(fix)
}

func fix(sd *pro.Speed) {
	fileInfos, err := ioutil.ReadDir(pro.CommonDir)
	if err != nil {
		log.Println(err)
	}
	for i := range fileInfos {
		if !fileInfos[i].IsDir() {
			if strings.HasSuffix(fileInfos[i].Name(), "fail_pic") {
				fixPic(fileInfos[i].Name(), sd)
			} else {
				fixWeb(fileInfos[i].Name(), sd)
			}
		}
	}
}

func fixPic(path string, sd *pro.Speed) {
	f, err := os.Open(pro.CommonDir + path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), "<->")
		img, dir := s[0], s[1]
		sd.Add(1)
		go pro.Download(img, dir, sd)
		time.Sleep(pro.Interval)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func fixWeb(path string, sd *pro.Speed) {
	f, err := os.Open(pro.CommonDir + path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		sd.WebAdd(1)
		go pro.Fetch(scanner.Text(), sd)
		time.Sleep(pro.Interval)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

}
