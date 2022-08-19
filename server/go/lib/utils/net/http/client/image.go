package client

import (
	"fmt"
	"github.com/actliboy/hoper/server/go/lib/utils/fs"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const DownloadKey = ".downloading"

func GetImage(url string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9;charset=utf-8")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36")
	var resp *http.Response
	for i := 0; i < 20; i++ {
		if i > 0 {
			time.Sleep(time.Second)
		}
		resp, err = defaultClient.Do(req)
		if err != nil {
			log.Println(err, "url:", url)
			continue
		}
		if resp.StatusCode != 200 {
			resp.Body.Close()
			return nil, fmt.Errorf("返回错误，状态码：%d,url:%s", resp.StatusCode, url)
		}
	}

	return resp.Body, nil
}

func DownloadImage(filepath, url string) error {
	reader, err := GetImage(url)
	if err != nil {
		return err
	}
	err = Download(filepath, reader)
	reader.Close()
	return err
}

func Download(filepath string, reader io.Reader) error {
	filepath = filepath + DownloadKey
	f, err := fs.Create(filepath)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, reader)
	if err != nil {
		f.Close()
		os.Remove(filepath)
		return err
	}

	err = f.Close()
	if err != nil {
		os.Remove(filepath)
		return err
	}
	return os.Rename(filepath, filepath[:len(filepath)-len(DownloadKey)])
}
