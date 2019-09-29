package main


/*
unsigned long long int fibonacci(unsigned int n);


unsigned long long int fibonacci(unsigned int n) {
 if (n < 2) return 1;
    return fibonacci(n - 1) + fibonacci(n - 2);
}

 */
import "C"
import (
	"fmt"
	"time"
)
//可能这就是cgo的应用场景吧，但是真的存在吗，n为34时，时间一致，n为42时，时间差了一倍
//go向c传参有额外开销
func main()  {
	useGo()
	useCGo()
}

func useGo()  {
	now:=time.Now()
	fmt.Println(now)
	fmt.Println(C.fibonacci(42))
	fmt.Println(time.Now().Sub(now))
}

func useCGo()  {
	now:=time.Now()
	fmt.Println(now)
	fmt.Println(goFib(42))
	fmt.Println(time.Now().Sub(now))
}

func goFib(n uint) uint64 {
	if n <2 {return 1}
	return goFib(n - 1) + goFib(n - 2)
}
