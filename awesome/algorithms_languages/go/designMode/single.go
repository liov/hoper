package main

import (
	"sync"
	"sync/atomic"
)

/*懒汉模式
懒汉模式是开源项目中使用最多的一种，最大的缺点是非线程安全的*/
type singleton struct {
}

// private
var instance *singleton

// public
func GetInstance() *singleton {
	if instance == nil {
		instance = &singleton{} // not thread safe
	}
	return instance
}

//带锁的单例模式

var mu sync.Mutex

func GetInstanceWithMutex() *singleton {
	mu.Lock()
	defer mu.Unlock()
	if instance == nil {
		instance = &singleton{} // unnecessary locking if
		//instance already created
	}
	return instance
}

//带检查锁的单例模式

var initialized uint32

func GetInstanceWithAtomic() *singleton {
	if atomic.LoadUint32(&initialized) == 1 {
		return instance
	}
	mu.Lock()
	defer mu.Unlock()
	if initialized == 0 {
		instance = &singleton{}
		atomic.StoreUint32(&initialized, 1)
	}
	return instance
}

//once
var once sync.Once

func GetInstanceOnce() *singleton {
	once.Do(func() {
		instance = &singleton{}
	})
	return instance
}
