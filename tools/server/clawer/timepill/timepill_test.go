package timepill

import (
	"fmt"
	"github.com/liov/hoper/server/go/lib/utils/encoding/json/iterator"
	"log"
	surl "net/url"
	"path"
	"strconv"
	"strings"
	"testing"
	"time"
	"tools/clawer/timepill/model"
)

func TestJson(t *testing.T) {
	var comments []*model.Comment
	err := iterator.Unmarshal([]byte(`[]`), &comments)
	fmt.Println(err, comments)
}

func TestGetDiaryComments(t *testing.T) {
	Token = "Basic bGJ5LmlAcXEuY29tOmxieTYwNA=="
	fmt.Println(ApiService.GetDiaryComments(6817247))
}

func TestTimeParse(f *testing.T) {
	fmt.Println(time.ParseInLocation("2006-01-02T15:04:05+08:00", "2010-03-18T13:03:48+08:00", time.Local))
}

func TestDownloadPic(t *testing.T) {
	Conf.TimePill.PhotoPath = "D:/F/timepill"
	DownloadPic(100774418, "http://s4.timepill.net/s/w640/photos/2022-08-24/b2s6hibqesmmq4opb8xbbgru8z4i3vwa.jpg", "2022-08-24 14:48:15")
}

func TestPicUrl(t *testing.T) {
	date := "2010-03-17"
	userId := 1
	notebookId := 2
	URL, _ := surl.Parse("http://s4.timepill.net/book_cover/0/15.jpg?v=1")
	v := URL.Query().Get("v")
	filepath := Conf.TimePill.PhotoPath
	originFileName := path.Base(URL.Path)
	originFileName = strings.TrimSuffix(originFileName, path.Ext(originFileName)) + "-v" + v + path.Ext(originFileName)
	filepath += strings.Join([]string{filepath, "book_cover", date[:4], date[:7], strconv.Itoa(userId)}, "/")
	filepath += "/" + strconv.Itoa(notebookId)
	filepath += "/" + originFileName
	log.Println(filepath)
}
