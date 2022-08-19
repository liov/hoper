package main

import (
	"github.com/jlaffaye/ftp"
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
