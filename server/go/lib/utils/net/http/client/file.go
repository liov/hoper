package client

import (
	"fmt"
	"github.com/liov/hoper/server/go/lib/utils/fs"
	httpi "github.com/liov/hoper/server/go/lib/utils/net/http"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const DownloadKey = ".downloading"

func GetFileWithReq(url string, setReq Option) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// 如果自己设置了接受编码，http库不会自动gzip解压，需要自己处理，不加Accept-Encoding和Range头会自动设置gzip
	//req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9;charset=utf-8")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36")

	if setReq != nil {
		setReq(req)
	}

	var resp *http.Response
	for i := 0; i < 3; i++ {
		if i > 0 {
			time.Sleep(time.Second)
		}
		resp, err = defaultClient.Do(req)
		if err != nil {
			log.Println(err, "url:", url)
			continue
		}
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			resp.Body.Close()
			return nil, fmt.Errorf("返回错误，状态码：%d,url:%s", resp.StatusCode, url)
		} else {
			break
		}
	}

	return resp.Body, nil
}

func DownloadImage(filepath, url string) error {
	return DownloadFileWithReq(filepath, url, ImageOption)
}

func DownloadFile(filepath, url string) error {
	return DownloadFileWithReq(filepath, url, nil)
}

func DownloadFileWithReq(filepath, url string, setReq Option) error {
	reader, err := GetFileWithReq(url, setReq)
	if err != nil {
		return err
	}
	err = Download(filepath, reader)
	reader.Close()
	return err
}

func DownloadFileWithRefer(filepath, url, refer string) error {
	return DownloadFileWithReq(filepath, url, SetRefer(refer))
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

func DownloadIfNotExists(filepath string, reader io.Reader) error {
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return Download(filepath, reader)
	}
	log.Println("已存在:", filepath)
	return err
}

func GetImage(url string) (io.ReadCloser, error) {
	return GetFileWithReq(url, ImageOption)
}

func GetFile(url string) (io.ReadCloser, error) {
	return GetFileWithReq(url, nil)
}

func GetFileWithRefer(url, refer string) (io.ReadCloser, error) {
	return GetFileWithReq(url, SetRefer(refer))
}

type Option func(req *http.Request)

func SetRefer(refer string) Option {
	return func(req *http.Request) {
		req.Header.Set(httpi.HeaderReferer, refer)
	}
}

func SetAccept(refer string) Option {
	return func(req *http.Request) {
		req.Header.Set(httpi.HeaderAccept, refer)
	}
}

func ImageOption(req *http.Request) {
	req.Header.Set(httpi.HeaderAccept, "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
}
