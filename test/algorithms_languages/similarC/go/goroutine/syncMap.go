package main

import (
	"fmt"
	"sync"
)

func main() {
	list := map[string]interface{}{
		"name":          "田馥甄",
		"birthday":      "1983年3月30日",
		"age":           34,
		"hobby":         []string{"听音乐", "看电影", "电视", "和姐妹一起讨论私人话题"},
		"constellation": "白羊座",
	}

	var m sync.Map
	for k, v := range list {
		m.Store(k, v)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		m.Store("age", 22)
		m.LoadOrStore("tag", 8888)
		wg.Done()
	}()

	go func() {
		m.Delete("constellation")
		m.Store("age", 18)
		wg.Done()
	}()

	wg.Wait()

	m.Range(func(key, value interface{}) bool {
		fmt.Println(key, value)
		return true
	})
}

//https://zhuanlan.zhihu.com/p/27642824
