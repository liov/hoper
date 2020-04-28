package pick

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/kataras/pio"
	"github.com/liov/hoper/go/v2/utils/net/http/api/apidoc"
	"github.com/liov/hoper/go/v2/utils/strings2"
)

type EasyRouter struct {
	mu           sync.RWMutex
	route        map[string][]*methodHandle
	es           []muxEntry
	middleware   Handlers
	NotFound     http.Handler
	PanicHandler func(http.ResponseWriter, *http.Request, interface{})
	hosts        bool
}

const MethodAny = "*"

type muxEntry struct {
	preUrl string
	handle []*methodHandle
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
				if _, h2 := getHandle(methodInfo.method, mh); h2.IsValid() {
					panic("url：" + methodInfo.path + "已注册")
				} else {
					mh = append(mh, &methodHandle{methodInfo.method, methodInfo.middleware, nil, value.Method(j)})
					router.route[methodInfo.path] = mh
				}
			} else {
				router.route[methodInfo.path] = []*methodHandle{{methodInfo.method, methodInfo.middleware, nil, value.Method(j)}}
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

	r.es = appendSorted(r.es, muxEntry{
		path,
		[]*methodHandle{{
			http.MethodGet,
			nil,
			http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				req.URL.Path = req.URL.Path[len(path):]
				fileServer.ServeHTTP(w, req)
			}),
			reflect.Value{},
		},
		},
	})
}

func (r *EasyRouter) Use(middleware ...http.Handler) {
	r.middleware = append(r.middleware, middleware...)
}

func (r *EasyRouter) Handle(method, path string, handle ...http.Handler) {
	newMh := &methodHandle{method, handle[:len(handle)-1], handle[len(handle)-1], reflect.Value{}}
	if mh, ok := r.route[path]; ok {
		if h, _ := getHandle(method, mh); h != nil {
			panic("url：" + path + "已注册")
		}
	} else {
		r.route[path] = append(mh, newMh)
		if path[len(path)-1] == '/' {
			r.es = appendSorted(r.es, muxEntry{path, []*methodHandle{newMh}})
		}
	}

	if path[0] != '/' {
		r.hosts = true
	}
	fmt.Printf(" %s\t %s %s\t %s\n",
		pio.Green("API:"),
		pio.Yellow(strings2.FormatLen(method, 6)),
		pio.Blue(strings2.FormatLen(path, 50)), pio.Purple(path))
}

func appendSorted(es []muxEntry, e muxEntry) []muxEntry {
	for i, mh := range es {
		if mh.preUrl == e.preUrl {
			es[i].handle = append(es[i].handle, e.handle...)
			return es
		}
	}

	n := len(es)
	i := sort.Search(n, func(i int) bool {
		return len(es[i].preUrl) < len(e.preUrl)
	})
	if i == n {
		return append(es, e)
	}
	// we now know that i points at where we want to insert
	es = append(es, muxEntry{}) // try to grow the slice in place, any entry works.
	copy(es[i+1:], es[i:])      // Move shorter entries down
	es[i] = e
	return es
}

func (r *EasyRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if r.PanicHandler != nil {
		defer r.recv(w, req)
	}
	r.middleware.ServeHTTP(w, req)
	if mh, ok := r.route[req.URL.Path]; ok {
		h1, h2 := getHandle(req.Method, mh)
		if h1 != nil {
			h1.ServeHTTP(w, req)
		}
		if h2.IsValid() {
			commonHandler(w, req, h2, nil)
		}
		return
	}
	for i := range r.es {
		if strings.HasPrefix(req.URL.Path, r.es[i].preUrl) {
			h1, _ := getHandle(req.Method, r.es[i].handle)
			h1.ServeHTTP(w, req)
			return
		}
	}

	if r.NotFound != nil {
		r.NotFound.ServeHTTP(w, req)
	} else {
		http.NotFound(w, req)
	}
}
