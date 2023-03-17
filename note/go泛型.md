1.一般不使用泛型扩展结构体
```go
type Foo[T any] struct{
	Prop T
}

type Bar = Foo[int]
```
而使用组合创建新类型
```go
type Foo struct{
	
}

type Bar struct{
    Foo
	Field int
}
```