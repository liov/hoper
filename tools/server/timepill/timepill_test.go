package timepill

import (
	"fmt"
	"github.com/actliboy/hoper/server/go/lib/utils/encoding/json"
	"testing"
	"time"
	"tools/timepill/model"
)

func TestJson(t *testing.T) {
	var comments []*model.Comment
	err := json.Unmarshal([]byte(`[]`), &comments)
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
