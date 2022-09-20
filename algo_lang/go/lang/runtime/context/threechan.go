package main

import (
	"context"
	"fmt"
	"time"
)

func A(ctx context.Context) int {
	ctx = context.WithValue(ctx, "AFunction", "Great")

	go B(ctx)

	select {
	// 监测自己上层的ctx ...
	case <-ctx.Done():
		fmt.Println("A Done")
		return -1
	}
	return 1
}

func B(ctx context.Context) int {
	fmt.Println("A value in B:", ctx.Value("AFunction"))
	ctx = context.WithValue(ctx, "BFunction", 999)

	go C(ctx)

	select {
	// 监测自己上层的ctx ...
	case <-ctx.Done():
		fmt.Println("B Done")
		return -2
	}
	return 2
}

func C(ctx context.Context) int {
	fmt.Println("B value in C:", ctx.Value("AFunction"))
	fmt.Println("B value in C:", ctx.Value("BFunction"))
	select {
	// 结束时候做点什么 ...
	case <-ctx.Done():
		fmt.Println("C Done")
		return -3
	}
	return 3
}

func main() {
	// 自动取消(定时取消)

	timeout := 10 * time.Second
	ctx, _ := context.WithTimeout(context.Background(), timeout)

	fmt.Println("A 执行完成，返回：", A(ctx))
	select {
	case <-ctx.Done():
		fmt.Println("context Done")
		break
	}

	time.Sleep(20 * time.Second)
}
