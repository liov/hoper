package pick

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/kataras/pio"
	"github.com/liov/hoper/go/v2/utils/net/http/api/apidoc"
	"github.com/liov/hoper/go/v2/utils/strings2"
)

type EasyRouter struct {
	route        map[string][]*methodHandle
	serveFiles   []serveFile
	middleware   HandlerFuncs
	NotFound     http.Handler
	PanicHandler func(http.ResponseWriter, *http.Request, interface{})
}

type serveFile struct {
	preUrl string
	handle http.HandlerFunc
}

func NewEasyRouter(genApi bool, modName string) *EasyRouter {
	router := &EasyRouter{
		route: make(map[string][]*methodHandle),
	}
	for _, v := range svcs {
		describe, preUrl, _ := v.Service()
		value := reflect.ValueOf(v)
		if value.Kind() != reflect.Ptr {
			log.Fatal("必须传入指针")
		}

		for j := 0; j < value.NumMethod(); j++ {
			method := value.Type().Method(j)
			if method.Type.NumIn() < 2 || method.Type.NumOut() != 2 {
				continue
			}
			methodInfo := getMethodInfo(value.Method(j))
			if methodInfo == nil {
				log.Fatalf("%s未注册", method.Name)
			}
			methodInfo.path, methodInfo.version = parseMethodName(method.Name)
			if methodInfo.version == 0 {
				methodInfo.version = 1
			}
			methodInfo.path = preUrl + "/" + methodInfo.path
			methodInfo.path = strings.Replace(methodInfo.path, "${version}", "v"+strconv.Itoa(methodInfo.version), 1)
			if methodInfo.path == "" || methodInfo.method == "" || methodInfo.title == "" || methodInfo.createlog.version == "" {
				log.Fatal("接口路径,方法,描述,创建日志均为必填")
			}
			if mh, ok := router.route[methodInfo.path]; ok {
				if getHandle(methodInfo.method, mh).IsValid() {
					panic("url：" + methodInfo.path + "已注册")
				} else {
					mh = append(mh, &methodHandle{methodInfo.method, value.Method(j)})
					router.route[methodInfo.path] = mh
				}
			} else {
				router.route[methodInfo.path] = []*methodHandle{{methodInfo.method, value.Method(j)}}
			}
			fmt.Printf(" %s\t %s %s\t %s\n",
				pio.Green("API:"),
				pio.Yellow(strings2.FormatLen(methodInfo.method, 6)),
				pio.Blue(strings2.FormatLen(methodInfo.path, 50)), pio.Purple(methodInfo.title))
			if genApi {
				methodInfo.Api(value.Method(j).Type(), describe, value.Type().Name())
				apidoc.WriteToFile(apidoc.FilePath, modName)
			}
		}
	}
	/*	if genApi {
		doc := doc(modName)
		router.route["/api-doc/md"], func(ctx iris.Context) {
			ctx.Text("[TOC]\n\n" + doc)
		})
	}*/
	registered()
	return router
}

func (r *EasyRouter) recv(w http.ResponseWriter, req *http.Request) {
	if rcv := recover(); rcv != nil {
		r.PanicHandler(w, req, rcv)
	}
}

func (r *EasyRouter) ServeFiles(path string, root string) {

	fileServer := http.FileServer(http.Dir(root))

	r.serveFiles = append(r.serveFiles, serveFile{
		path,
		func(w http.ResponseWriter, req *http.Request) {
			req.URL.Path = req.URL.Path[len(path):]
			fileServer.ServeHTTP(w, req)
		},
	})
}

func (r *EasyRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if r.PanicHandler != nil {
		defer r.recv(w, req)
	}
	r.middleware.ServeHTTP(w, req)
	if mh, ok := r.route[req.URL.Path]; ok {
		if handle := getHandle(req.Method, mh); handle.IsValid() {
			commonHandler(w, req, handle, nil)
			return
		}
	}
	for i := range r.serveFiles {
		if strings.HasPrefix(req.URL.Path, r.serveFiles[i].preUrl) {
			r.serveFiles[i].handle(w, req)
			return
		}
	}

	if r.NotFound != nil {
		r.NotFound.ServeHTTP(w, req)
	} else {
		http.NotFound(w, req)
	}
}
