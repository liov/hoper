go 的抽象能力所限制
任务调度框架的task如果想携带额外的属性，task的定义方式

# interface{}字段
```go
type TaskMeta struct {
    // [field: type]
}

type TaskFunc func(ctx context.Context) ([]*Task, error)

type Task struct {
	TaskMeta
	TaskFunc
	Props interface{}
}

```
简单
类型不直观，要使用需要断言

# interface实现
```go
type TaskInterface interface {
	HasTask() *Task
}

type TaskFunc func(ctx context.Context) ([]TaskInterface, error)

type Task struct {
    TaskMeta
    TaskFunc
}

type Custom struct{
	Prop string
}

func (c *Custom) HasTask() *Task {
return &Task{}
}

```
框架承担处理interface的成本,使用方可以清楚看到task其他属性类型
不安全，框架要对HasTask可能返回为空处理，使用方要注意返回的是接口切片

# 泛型
```go
type Task[T any] struct {
    TaskMeta
    TaskFunc
	Prop T
}
```
虽然go的泛型很弱，目前看来还是优于以上两种的