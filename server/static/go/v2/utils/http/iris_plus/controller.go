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
		value1.Field(0).Set(reflect.ValueOf(handler))

		for j := 0; j < value.NumMethod(); j++ {
			if value.Type().Method(j).Name == "Middle" {
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
	Middle []iris.Handler
	app    *iris.Application
}

type apiInfo struct {
	path, method, describe, auth, date string
	//生成api用的参数，不优雅的实现
	request, response, service interface{}
	version                    int
	changelog                  []changelog
	createlog                  changelog
}

type changelog struct {
	version, date, log string
}

func (h *Handler) Path(p string) *Handler {
	h.path = p
	return h
}

func (h *Handler) Date(d string) *Handler {
	h.date = d
	return h
}

func (h *Handler) ChangeLog(v, date, log string) *Handler {
	h.changelog = append(h.changelog, changelog{v, date, log})
	return h
}

func (h *Handler) CreateLog(v, date, log string) *Handler {
	h.createlog = changelog{v, date, log}
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

func (h *Handler) Auth(a string) *Handler {

	h.auth = a
	return h
}

func (h *Handler) Handle(hs ...iris.Handler) *Handler {
	path := "/api/v" + strconv.Itoa(h.version) + h.path
	handles := append(h.Middle, hs...)
	h.app.Handle(h.method, path, handles...)
	fmt.Printf(" %s\t %s %s\t %s\n",
		pio.Purple("API:"),
		pio.Yellow(strings2.FormatLen(h.apiInfo.method, 6)),
		pio.Blue(path), pio.Gray(h.apiInfo.describe))
	return h
}

func (h *Handler) Api() *Handler {
	return h
}
