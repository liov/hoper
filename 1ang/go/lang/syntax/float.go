package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Printf("%.2f\n", 199.995)               // 200.00
	fmt.Printf("%.2f\n", 11198.995)             // 11199.00
	fmt.Printf("%.2f\n", 21198.995)             // 21198.99
	fmt.Printf("%.0f\n", 26.5)                  // 26
	fmt.Printf("%.0f\n", 27.5)                  // 26
	fmt.Printf("%.1f\n", 2.65)                  // 2.6
	fmt.Printf("%.2f\n", 0.265)                 // 0.27
	fmt.Printf("%.0f\n", 26.6)                  // 27
	fmt.Printf("%.0f\n", 26.51)                 // 0.27
	fmt.Printf("%.0f\n", 26.5+0.01)             // 0.27
	fmt.Printf("%.0f\n", 27.5+0.01)             // 0.27
	fmt.Printf("%.0f\n", 27*(1/3)+0.01)         // 0.27
	fmt.Println(math.Round(205.49999999999999)) // 0.27
	fmt.Println(206-205.49999999999999 > 0.5)
	fmt.Println(float64(3) / float64(1384) * float64(346))
	fmt.Println(float64(9) / float64(4152) * float64(1038))
	fmt.Println(float64(18) / float64(8304) * float64(2076))
	fmt.Println(math.Round(float64(3) / float64(1384) * float64(346)))
	fmt.Println(math.Round(float64(9) / float64(4152) * float64(1038)))
	fmt.Println(math.Round(float64(18) / float64(8304) * float64(2076))) // 可以准确算出4.5
	/*for i := 7; i < 10000; i++ {
		for j := i; j < 10000; j++ {
			for x := 15; x < j; x++ {
				num := float64(i) / float64(j) * float64(x)
				diff := num - math.Floor(num)
				if diff > 0.49999 && diff < 0.5 {
					fmt.Println(i, j, x, num, math.Round(num))
				}
			}
		}
	}*/

}
