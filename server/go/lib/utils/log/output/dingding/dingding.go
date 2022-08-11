package dingding

import (
	"fmt"
	"github.com/actliboy/hoper/server/go/lib/utils/encoding/json"
	"github.com/actliboy/hoper/server/go/lib/utils/strings"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/url"
	"time"
)

type TalkHook struct {
	AcceptedLevels []zapcore.Level
	BotURL         string
	AtPhoneArr     []string //填写在atMobiles中的@ 手机号
	AtPhoneStr     string   //填写在text中的@ 手机号
}

func (th *TalkHook) Write(b []byte) (n int, err error) {
	return
}

func (th *TalkHook) Sync() error {
	return nil
}

func (th *TalkHook) Close() error {
	return nil
}

type KibanaField struct {
	Id          string `json:"id"`
	Interface   string `json:"interface"`
	Method      string `json:"method"`
	Param       string `json:"param"`
	Other       string `json:"other"`
	ProcessTime string `json:"processTime"`
	Result      string `json:"result"`
}

func (th *TalkHook) FormatLog(e *zapcore.Entry) error {

	return nil
}

// dingding://${token}?sercret=${sercret}
func RegisterSink() {
	_ = zap.RegisterSink("dingding", func(url *url.URL) (sink zap.Sink, e error) {
		th := new(TalkHook)
		return th, nil
	})
}

func NewCore() zapcore.Core { return core{} }

type core struct {
}

func (core) Enabled(zapcore.Level) bool          { return false }
func (n core) With([]zapcore.Field) zapcore.Core { return n }
func (core) Check(e zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	fields := new(KibanaField)
	json.Unmarshal(stringsi.ToBytes(e.Message), fields)
	t := e.Time
	url := fmt.Sprintf(viper.GetString("domain.kibana")+"/app/kibana#/discover?_g=%28refreshInterval:%28display:Off,pause:!f,value:0%29,time:%28from:'%s',mode:absolute,to:'%s'%29%29&_a=%28columns:!%28_source%29,index:'8c289dc0-83ed-11e8-9c6f-f5bb9ff824a1',interval:auto,query:%28language:lucene,query:'@source:%%22%+v%%22%%20%%20AND%%20@fields.uuid:%%22[%+v]%%22'%29%29",
		t.Add(-time.Minute).Format("2006-01-02T15:04:05+08:00"), t.Add(time.Minute).Format("2006-01-02T15:04:05+08:00"), e.Caller.String(), fields.Id,
	)

	fmt.Sprintf(`
**项目:** %+v

**时间:** %s

**接口:** %+v

**错误:** %+v

**环境:** %s

  
  
[点击查看](%s)
`, e.Caller.String(), t.Format("2006-01-02 15:04"), fields.Interface,
		fields.Result, viper.GetString("env"), url)

	//if f, ok := e.Data["@fields"].(map[string]interface{}); ok {
	//	if stack, ok := f["callstack"].(string); ok {
	//		info +=  "\n\n\n```go\n" + stack + "\n```"
	//	}
	//}
	return ce
}
func (core) Write(zapcore.Entry, []zapcore.Field) error { return nil }
func (core) Sync() error                                { return nil }