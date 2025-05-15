package main

import (
	"errors"
	"github.com/bdragon300/tusgo"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

func UploadWithRetry(dst *tusgo.UploadStream, src *os.File) error {
	// Set stream and file pointer to be equal to the remote pointer
	// (if we resume the upload that was interrupted earlier)
	if _, err := dst.Sync(); err != nil {
		return err
	}
	if _, err := src.Seek(dst.Tell(), io.SeekStart); err != nil {
		return err
	}

	_, err := io.Copy(dst, src)
	attempts := 10
	for err != nil && attempts > 0 {
		if _, ok := err.(net.Error); !ok && !errors.Is(err, tusgo.ErrChecksumMismatch) {
			return err // Permanent error, no luck
		}
		time.Sleep(5 * time.Second)
		attempts--
		_, err = io.Copy(dst, src) // Try to resume the transfer again
	}
	if attempts == 0 {
		return errors.New("too many attempts to upload the data")
	}
	return nil
}

func CreateUploadFromFile(f *os.File, cl *tusgo.Client) *tusgo.Upload {
	finfo, err := f.Stat()
	if err != nil {
		panic(err)
	}

	u := tusgo.Upload{}
	if _, err = cl.CreateUpload(&u, finfo.Size(), true, map[string]string{"filename": "zls-x86_64-windows.zip"}); err != nil {
		panic(err)
	}
	return &u
}

func main() {
	baseURL, _ := url.Parse("http://localhost:8080/files/")
	cl := tusgo.NewClient(http.DefaultClient, baseURL)
	cl.GetRequest = func(method, url string, body io.Reader, tusClient *tusgo.Client, httpClient *http.Client) (*http.Request, error) {
		req, err := http.NewRequest(method, url, body)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "xxx")
		return req, nil
	}
	f, err := os.Open(`D:\Download\zls-x86_64-windows.zip`)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	u := CreateUploadFromFile(f, cl)

	stream := tusgo.NewUploadStream(cl, u)
	if err = UploadWithRetry(stream, f); err != nil {
		panic(err)
	}
}
