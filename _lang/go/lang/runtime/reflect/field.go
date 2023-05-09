package main

import (
	"fmt"
	"reflect"
)

type FlagConfig struct {
	p *PreSet
	a AfterSet
}

type PreSet struct {
	// environment
	Env string
	// 配置文件路径
	ConfUrl string
	// 是否监听配置文件
	Watch bool
	// 代理, socks5://localhost:1080
	Proxy string
}

type AfterSet struct {
	Watch bool `FlagConfig:"name:awatch"`
}

type FlagConfig2 struct {
	p *PreSet2
	a AfterSet2
}

type PreSet2 struct {
	// environment
	Env string
	// 配置文件路径
	ConfUrl string
	// 是否监听配置文件
	Watch2 bool
	// 代理, socks5://localhost:1080
	Proxy2 string
}

type AfterSet2 struct {
	Watch2 bool `FlagConfig:"name:awatch"`
}

func main() {
	flag := FlagConfig{}
	vaule := reflect.ValueOf(&flag).Elem()

	for i := 0; i < vaule.NumField(); i++ {
		fmt.Println(vaule.Field(i).Type())
	}

	flag1 := FlagConfig{
		p: &PreSet{Env: "dev"},
		a: AfterSet{Watch: true},
	}
	flag2 := FlagConfig2{}
	reflect.ValueOf(&flag2).Elem().Set(reflect.ValueOf(&flag1).Elem())
	flag1.p.Env = "prod"
	flag1.a.Watch = false
	fmt.Println(flag1)
	fmt.Println(flag2)

	flag3 := FlagConfig{
		p: &PreSet{Env: "dev"},
		a: AfterSet{Watch: true},
	}
	flag4 := FlagConfig2{}
	reflect.ValueOf(&flag4).Elem().Set(reflect.ValueOf(flag3))
	flag3.p.Env = "prod"
	flag3.a.Watch = false
	fmt.Println(flag3)
	fmt.Println(flag4)
}
