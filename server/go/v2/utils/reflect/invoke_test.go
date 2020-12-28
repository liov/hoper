package reflecti

import (
	"fmt"
	"testing"

	"github.com/liov/hoper/go/v2/utils/encoding/json"
)

type InvokeFoo struct{}
type InvokeBar struct{}

func (b *InvokeBar) BarFuncAdd(argOne, argTwo float64) float64 {

	return argOne + argTwo
}

func (f *InvokeFoo) FooFuncSwap(argOne, argTwo string) (string, string) {

	return argTwo, argOne
}

func TestInvokeByValues(t *testing.T) {
	foo := &InvokeFoo{}
	bar := &InvokeBar{}
	reflectinvoker := NewReflectinvoker()
	reflectinvoker.RegisterMethod(foo)
	reflectinvoker.RegisterMethod(bar)
	req := Request{FuncName: "InvokeFoo.FooFuncSwap", Params: []interface{}{"1", "2"}}
	data, _ := json.Standard.Marshal(req)
	resultJson := reflectinvoker.InvokeByJson(data)
	result := Response{}
	err := json.Standard.Unmarshal(resultJson, &result)
	if err != nil {
		t.Log(err)
	}
	info := "handleJsonMessage get a result\n"
	info += "raw:\n" + string(resultJson) + "\n"
	info += "function: " + result.FuncName + " \n"
	info += fmt.Sprintf("result: %v\n", result.Result)
	info += fmt.Sprintf("error: %s\n", result.ErrorMsg)

	t.Log(info)
}
