package main

import "fmt"

// 变参函数，参数不定长
func list(nums ...int) {
	nums[0] = 0
	fmt.Println(nums)
}

func main() {
	// 常规调用，参数可以多个
	list(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	// 在参数同类型时，可以组成slice使用 parms... 进行参数传递
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	list(numbers...) // slice时使用
	//变参函数传的也是引用
	fmt.Println(numbers)

	x := [100000]int{1, 2, 3}

	func(arr [100000]int) {
		arr[0] = 7
		fmt.Println(arr) //prints [7 2 3]
	}(x)

	fmt.Println(x) //prints [1 2 3] (not ok if you need [7 2 3])
}
