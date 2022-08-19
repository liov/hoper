package fs

import "os"

type DirEntities []os.DirEntry

func (e DirEntities) Len() int {
	return len(e)
}

func (e DirEntities) Less(i, j int) bool {
	filei, _ := e[i].Info()
	filej, _ := e[j].Info()
	return filei.ModTime().After(filej.ModTime())
}

func (e DirEntities) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}
