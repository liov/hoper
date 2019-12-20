package iris_plus

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/kataras/pio"
	"github.com/liov/hoper/go/v2/utils/strings2"
)

func ApiDoc(ctrl []Controller) {
	handler := &Handler{apiInfo: &apiInfo{}}
	for i := range ctrl {
		value := reflect.ValueOf(ctrl[i])
		if value.Kind() != reflect.Ptr {
			panic("必须传入指针")
		}

		value1 := value.Elem()
		if value1.NumField() == 0 {
			panic("controller必须有一个类型为Handler的field")
		}
		if value1.Field(0).Type() != reflect.TypeOf(handler) {
			panic("Handler field必须在第一个")
		}
		value1.Field(0).Set(reflect.ValueOf(handler))

		for j := 0; j < value.NumMethod(); j++ {
			value.Method(j).Call(nil)
		}

		path := "/api/v" + strconv.Itoa(handler.version) + handler.path
		fmt.Println(path)
	}
}

func Register(app *iris.Application, ctrl []Controller) {
	handler := &Handler{apiInfo: &apiInfo{}, app: app}
	for i := range ctrl {
		value := reflect.ValueOf(ctrl[i])
		if value.Kind() != reflect.Ptr {
			panic("必须传入指针")
		}

		value1 := value.Elem()
		if value1.NumField() == 0 {
			panic("controller必须有一个类型为Handler的field")
		}
		if value1.Field(0).Type() != reflect.TypeOf(handler) {
			panic("Handler field必须在第一个")
		}
		handler.middle = ctrl[i].Middle()
		value1.Field(0).Set(reflect.ValueOf(handler))

		for j := 0; j < value.NumMethod(); j++ {
			method := value.Type().Method(j)
			if method.Name == "Middle" {
				continue
			}
			if method.Type.NumOut() == 1 && method.Type.Out(0) == reflect.TypeOf(handler) {
				continue
			}
			value.Method(j).Call(nil)
		}
	}
	handler = nil
}

type Controller interface {
	Middle() []iris.Handler
}

type Handler struct {
	*apiInfo
	middle []iris.Handler
	app    *iris.Application
}

type apiInfo struct {
	path, method, describe, date string
	//生成api用的参数，不优雅的实现
	request, response, service interface{}
	version                    int
	changelog                  []changelog
	createlog                  changelog
}

type changelog struct {
	version, auth, date, log string
}

func (h *Handler) Path(p string) *Handler {
	h.path = p
	return h
}

func (h *Handler) Date(d string) *Handler {
	h.date = d
	return h
}

func (h *Handler) ChangeLog(v, auth, date, log string) *Handler {
	h.changelog = append(h.changelog, changelog{v, auth, date, log})
	return h
}

func (h *Handler) CreateLog(v, auth, date, log string) *Handler {
	if h.createlog.version != "" {
		panic("创建记录只允许一条")
	}
	h.createlog = changelog{v, auth, date, log}
	return h
}

func (h *Handler) Version(v int) *Handler {
	h.version = v
	return h
}

func (h *Handler) Method(m string) *Handler {
	h.method = m
	return h
}

func (h *Handler) Describe(d string) *Handler {
	h.describe = d
	return h
}

func (h *Handler) Request(r interface{}) *Handler {
	h.request = r
	return h
}

func (h *Handler) Response(r interface{}) *Handler {
	h.response = r
	return h
}

func (h *Handler) Service(s interface{}) *Handler {
	h.service = s
	return h
}

func (h *Handler) Handle(hs ...iris.Handler) *Handler {
	if h.app == nil {
		return h
	}
	path := "/api/v" + strconv.Itoa(h.version) + h.path
	handles := append(h.middle, hs...)
	h.app.Handle(h.method, path, handles...)
	fmt.Printf(" %s\t %s %s\t %s\n",
		pio.Green("API:"),
		pio.Yellow(strings2.FormatLen(h.method, 6)),
		pio.Blue(strings2.FormatLen(path, 50)), pio.Purple(h.describe))
	return h
}

func (h *Handler) Api() *Handler {
	return h
}

type HandlerFunc func(*Handler)

func (handler HandlerFunc) Handle(hs ...iris.Handler) HandlerFunc {
	return func(h *Handler) {
		handler(h)
		if h.app == nil {
			return
		}
		path := "/api/v" + strconv.Itoa(h.version) + h.path
		handles := append(h.middle, hs...)
		h.app.Handle(h.method, path, handles...)
		fmt.Printf(" %s\t %s %s\t %s\n",
			pio.Purple("API:"),
			pio.Yellow(strings2.FormatLen(h.method, 6)),
			pio.Blue(strings2.FormatLen(path, 50)), pio.Green(h.describe))
	}
}
