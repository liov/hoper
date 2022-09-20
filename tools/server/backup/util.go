package backup

import (
	"github.com/actliboy/hoper/server/go/lib/utils/fs"
	"github.com/jlaffaye/ftp"
	"os"
)

type Entities []*ftp.Entry

func (e Entities) Len() int {
	return len(e)
}

func (e Entities) Less(i, j int) bool {
	return e[i].Time.After(e[j].Time)
}

func (e Entities) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func LastFile(dir string) (os.FileInfo, map[string]os.FileInfo, error) {
	entities, err := os.ReadDir(dir)
	if len(entities) == 0 {
		return nil, nil, err
	}
	var max int
	var maxDate string
	for _, entity := range entities {
		if !entity.IsDir() {
			continue
		}
		date := entity.Name()
		var dateNum int
		for _, c := range date {
			dateNum = dateNum*10 + int(c-'0')
		}
		if dateNum > max {
			max = dateNum
			maxDate = date
		}
	}
	return fs.LastFile(dir + "\\" + maxDate)
}
