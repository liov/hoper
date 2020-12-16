package main

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/liov/hoper/go/v2/utils/fs"
	py "tools/pinyin"
	"tools/pro"
)

type Result struct {
	TId         int
	Auth, Title string
	PicNum      int
}

func main() {
	page := 0
	pageSize := 1000
	wg := new(sync.WaitGroup)
	chanArray := createFile(wg)
	pro.SetDB()
	for {
		var results []*Result
		pro.DB.Table(`post`).Select(`t_id,auth,title,pic_num`).
			Where(`status = 0 AND pic_num > 0`).Order(`t_id`).
			Limit(pageSize).Offset(page * pageSize).Scan(&results)
		for _, result := range results {
			if result.PicNum > 60 {
				result.PicNum = 60
			}
			chanArray[result.PicNum/10] <- fixPath(result)
		}
		if len(results) < pageSize {
			break
		}
		page++
	}
	for i := range chanArray {
		close(chanArray[i])
	}
	wg.Wait()
}
func createFile(wg *sync.WaitGroup) []chan string {
	var chanArray = make([]chan string, 7)
	for i := range chanArray {
		chanArray[i] = make(chan string, 20)
	}
	wg.Add(len(chanArray))
	for i := range chanArray {
		go func(i int) {
			path := pro.CommonDir + strconv.Itoa(i) + pro.Ext
			f, _ := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
			for txt := range chanArray[i] {
				log.Println(txt)
				f.WriteString(txt + "\n")
			}
			wg.Done()
			f.Close()
		}(i)
	}
	return chanArray
}

func fixPath(result *Result) string {
	auth := pro.FixPath(result.Auth)
	title := pro.FixPath(result.Title)
	var num string
	switch {
	case result.TId >= 400000:
		num = "4"
	case result.TId >= 300000 && result.TId < 400000:
		num = "3"
	case result.TId < 100000:
		num = "1"
	default:
		num = "2"
	}
	dir := `F:\pic_` + num + pro.Sep + py.FistLetter(auth) + pro.Sep + auth + pro.Sep + title + `_` + strconv.Itoa(result.TId) + pro.Sep
	return fs.PathClean(dir)
}
