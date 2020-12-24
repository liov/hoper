package pick

import (
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/strings2"
)

type apiInfo struct {
	path, method, title string
	version             int
	changelog           []changelog
	createlog           changelog
	deprecated          *changelog
	middleware          []http.HandlerFunc
}

type changelog struct {
	version, auth, date, log string
}

func Path(p string) *apiInfo {
	return &apiInfo{path: p}
}

func Method(m string) *apiInfo {
	return &apiInfo{method: m}
}

func (api *apiInfo) Method(m string) *apiInfo {
	api.method = m
	return api
}

func (api *apiInfo) ChangeLog(v, auth, date, log string) *apiInfo {
	api.changelog = append(api.changelog, changelog{v, auth, date, log})
	return api
}

func (api *apiInfo) CreateLog(v, auth, date, log string) *apiInfo {
	if api.createlog.version != "" {
		panic("创建记录只允许一条")
	}
	api.createlog = changelog{v, auth, date, log}
	return api
}

func (api *apiInfo) Title(d string) *apiInfo {
	api.title = d
	return api
}

func (api *apiInfo) Version(v int) *apiInfo {
	api.version = v
	return api
}

func (api *apiInfo) Deprecated(v, auth, date, log string) *apiInfo {
	api.deprecated = &changelog{v, auth, date, log}
	return api
}

func (api *apiInfo) Middleware(m ...http.HandlerFunc) *apiInfo {
	api.middleware = m
	return api
}

//获取负责人
func (api *apiInfo) getPrincipal() string {
	if len(api.changelog) == 0 {
		return api.createlog.auth
	}
	if api.deprecated != nil {
		return api.deprecated.auth
	}
	return api.changelog[len(api.changelog)-1].auth
}


//简直就是精髓所在，真的是脑洞大开才能想到
func getMethodInfo(fv reflect.Value,preUrl string) (info *apiInfo) {
	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(*apiInfo); ok {
				info = v
				if info.version == 0 {
					info.version = 1
				}
				_, info.version = parseMethodName(fv.Type().Name())
				info.path = preUrl + "/v" + strconv.Itoa(info.version) + info.path
			} else {
				log.Panic(err)
			}
		}
	}()
	methodType := fv.Type()
	params := make([]reflect.Value, 0, fv.Type().NumIn())
	numIn := methodType.NumIn()
	numOut := methodType.NumOut()
	if numIn == 1 {
		panic("method至少一个参数且参数必须实现Claims接口")
	}
	if numIn > 2 {
		panic("method参数最多为两个")
	}
	if numOut != 2 {
		panic("method返回值必须为两个")
	}
	if !methodType.In(0).Implements(claimsType) {
		panic("service第一个参数必须实现Claims接口")
	}
	if !methodType.Out(1).Implements(errorType) {
		panic("service第二个返回值必须为error类型")
	}
	for i := 0; i < numIn; i++ {
		params = append(params, reflect.New(methodType.In(i).Elem()))
	}
	fv.Call(params)
	return nil
}

// 从方法名称分析出接口名和版本号
func parseMethodName(originName string) (name string, version int) {
	idx := strings.LastIndexByte(originName, 'V')
	version = 1
	if idx > 0 {
		if v, err := strconv.Atoi(originName[idx+1:]); err == nil {
			version = v
		} else {
			idx = len(originName)
		}
	} else {
		idx = len(originName)
	}
	name = strings2.LowerFirst(originName[:idx])
	return
}

