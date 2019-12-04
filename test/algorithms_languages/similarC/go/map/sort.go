//+build go1.12

package main

import "fmt"

func main() {
	//go1.12会顺序打印map
	fmt.Println(map[int]int{1: 1, 5: 2, 3: 3, 2: 2, 999: 999, 6: 6})
	//''是rune
	fmt.Println(map[byte]int{'z': 1, 'f': 2, 's': 3, 2: 'e', 'a': 999, 'p': 6})

	fmt.Println(map[rune]int{'z': 1, 'f': 2, 's': 3, 2: 'e', 'a': 999, 'p': 6})

	fmt.Println(map[int]int{'z': 1, 'f': 2, 's': 3, 2: 'e', 'a': 999, 'p': 6})

	fmt.Println(map[string]int{"z": 1, "f": 2, "s": 3, "e": 2, "a": 999, "p": 6})
	//遍历随机
	for k, v := range map[string]int{"z": 1, "f": 2, "s": 3, "e": 2, "a": 999, "p": 6} {
		fmt.Println(k, v)
	}

	fmt.Println('a' + 3)
}
