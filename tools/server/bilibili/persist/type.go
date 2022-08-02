package persist

import (
	"fmt"
	"sync"
	"tools/bilibili/tool"
)

type Item struct {
	Payload interface{}
}

type GetItemChan func(*sync.WaitGroup) (chan *Item, error)

func GetItemProcessFun() GetItemChan {
	var itemProcessFun GetItemChan
	if !tool.CheckFfmegStatus() {
		fmt.Println("Can't locate your ffmeg.The video your download can't be merged")
		itemProcessFun = VideoItemCleaner
	} else {
		itemProcessFun = VideoItemProcessor
	}

	return itemProcessFun
}
