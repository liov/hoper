package tnsq

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/nsqio/go-nsq"
)

var producer *nsq.Producer

/*// 主函数
func NsqpSend() {
	//读取控制台输入
	reader := bufio.NewReader(os.Stdin)
	data, _, _ := reader.ReadLine()
	command := string(data)
}
*/
// 初始化生产者
/*func init() {
	var err error
	producer, err = nsq.NewProducer(Addr4150, nsq.NewConfig())
	if err != nil {
		panic(err)
	}
	foo := &Foo{}
	bar := &Bar{}

	reflectinvoker = reflectinvoke.NewReflectinvoker()
	reflectinvoker.RegisterMethod(foo)
	reflectinvoker.RegisterMethod(bar)
	NewConsumer("topic_json","jchan", handleJsonMessage)
	NewConsumer("topic_string","schan", handleStringMessage)
}*/

//发布消息
func publish(topic string, message string) error {
	var err error
	if producer != nil {
		if len(message) == 0 { //不能发布空串，否则会导致error
			return nil
		}
		err = producer.Publish(topic, []byte(message)) // 发布消息
		return err
	}
	return fmt.Errorf("producer is nil", err)
}

func Start(c iris.Context) {
	stringType := c.URLParam("st")
	if stringType == "0" {
		message := c.URLParam("message")
		publish("topic_string", message)
	} else {
		message := c.FormValue("message")
		publish("topic_json", message)
	}
}
