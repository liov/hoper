package nsq

import (
	"fmt"

	"github.com/liov/hoper/go/v2/utils/json"
	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/reflect3"
	"github.com/nsqio/go-nsq"
)

func handleJsonMessage(message *nsq.Message) error {
	reflectinvoker := reflect3.NewReflectinvoker()
	resultJson := reflectinvoker.InvokeByJson(message.Body)
	result := reflect3.Response{}
	err := json.Json.Unmarshal(resultJson, &result)
	if err != nil {
		return err
	}
	info := "handleJsonMessage get a result\n"
	info += "raw:\n" + string(resultJson) + "\n"
	info += "function: " + result.FuncName + " \n"
	info += fmt.Sprintf("result: %v\n", result.Result)
	info += fmt.Sprintf("error: %s\n", result.ErrorMsg)

	log.Info(info)

	return nil
}
