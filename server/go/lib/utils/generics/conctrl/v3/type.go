package v3

import "github.com/actliboy/hoper/server/go/lib/utils/conctrl"

type Task interface {
	GetId() uint
	GetKind() conctrl.Kind
	NewTask() *conctrl.Task
}
