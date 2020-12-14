package test

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

func TestDelDir(t *testing.T) {
	fileInfos, err := ioutil.ReadDir(`F:\pic_2\★\★★★★\`)
	if err != nil {
		log.Println(err)
	}
	for i := range fileInfos {
		log.Println(fileInfos[i].Name())
		if strings.HasPrefix(fileInfos[i].Name(), `如懿传系列`) {
			err := os.RemoveAll(`F:\pic_2\★\★★★★\` + fileInfos[i].Name())
			if err != nil {
				t.Log(err)
			}
		}
	}
}
