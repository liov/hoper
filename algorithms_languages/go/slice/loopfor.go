package main

import "fmt"

func main() {
	m := map[int][]int{}
	//下级
	var subordinates = []int{1}
	//标志位
	var index = 0
	for {
		if len(subordinates[index:]) == 0 {
			break
		}
		for _, id := range subordinates[index:] {
			if ids, ok := m[id]; ok {
				subordinates = append(subordinates, ids...)
				index += len(ids)
			}
		}
	}
	fmt.Println(subordinates)
}
