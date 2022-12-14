package fs

import (
	"errors"
	"io"
	"net/http"
	"path"
	"strconv"
	"time"
)

func FetchFile(r *http.Request) (*FileInfo, error) {
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	vbytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	var file FileInfo
	file.Binary = vbytes
	file.name = path.Base(resp.Request.URL.Path)
	file.modTime, _ = time.ParseInLocation(time.RFC1123, resp.Header.Get("Last-Modified"), time.Local)
	file.size, _ = strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	return &file, nil
}
