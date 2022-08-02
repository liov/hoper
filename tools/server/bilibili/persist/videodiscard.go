package persist

import (
	"log"

	"sync"
)

func VideoItemCleaner(wgOutside *sync.WaitGroup) (chan *Item, error) {
	out := make(chan *Item)
	go func() {
		defer wgOutside.Done()
		itemCount := 0
		for item := range out {
			log.Printf("Item Saver:got item "+
				"#%d: %v", itemCount, item)
			itemCount++
		}
	}()
	return out, nil
}
