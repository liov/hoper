package client

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

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
	var reader io.ReadCloser
	var resp *http.Response
	for i := 0; i < 20; i++ {
		if i > 0 {
			time.Sleep(time.Second)
		}
		resp, err = client.Do(req)
		if err != nil {
			log.Println(err, "url:", url)
			continue
		}
		if resp.StatusCode != 200 {
			resp.Body.Close()
			return nil, fmt.Errorf("返回错误，状态码：%d,url:%s", resp.StatusCode, url)
		}

		if resp.Header.Get("Content-Encoding") == "gzip" {
			reader, err = gzip.NewReader(resp.Body)
			if err != nil {
				if resp != nil {
					resp.Body.Close()
				}
				log.Println(err, "url:", url)
				continue
			}
		} else {
			reader = resp.Body
		}
		if reader != nil {
			break
		}
	}
	if reader == nil {
		if resp != nil {
			resp.Body.Close()
		}
		msg := "请求失败：" + url
		if err != nil {
			msg = err.Error() + msg
		}
		return nil, errors.New(msg)
	}
	return reader, nil
}

func DownloadImage(filepath, url string) error {
	reader, err := GetImage(url)
	if err != nil {
		return err
	}
	defer reader.Close()
	dir := path.Dir(filepath)
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0666)
		if err != nil {
			return err
		}
	}
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, reader)
	if err != nil {
		return err
	}
	return nil
}
