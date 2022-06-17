package fs

import (
	"io"
	"io/fs"
	"net/http"
	"time"
)

type File struct {
	File http.File
	Name string
}

type FileInterface interface {
	io.Reader
	Name() string
}

type FileInfo struct {
	name    string
	modTime time.Time
	size    int64
	mode    fs.FileMode
	Binary  []byte
}

func (f *FileInfo) Name() string {
	return f.name
}

func (f *FileInfo) Size() int64 {
	return f.size
}

func (f *FileInfo) Mode() fs.FileMode {
	return f.mode
}

func (f *FileInfo) ModTime() time.Time {
	return f.modTime
}

func (f *FileInfo) IsDir() bool {
	return false
}

func (f *FileInfo) Sys() any {
	return nil
}
