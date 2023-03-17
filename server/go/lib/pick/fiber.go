package pick

import (
	"context"
	"encoding/json"
	"github.com/liov/hoper/server/go/lib/context/fasthttp_context"
	http_fs "github.com/liov/hoper/server/go/lib/utils/net/http/fs"
	"io"
	"net/http"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"github.com/liov/hoper/server/go/lib/protobuf/errorcode"
	"github.com/liov/hoper/server/go/lib/utils/log"
	httpi "github.com/liov/hoper/server/go/lib/utils/net/http"
	"github.com/liov/hoper/server/go/lib/utils/net/http/api/apidoc"
	fiber_build "github.com/liov/hoper/server/go/lib/utils/net/http/fasthttp/fiber"
)

type FiberService interface {
	//返回描述，url的前缀，中间件
	FiberService() (describe, prefix string, middleware []fiber.Handler)
}

var fiberSvcs = make([]FiberService, 0)

func RegisterFiberService(svc ...FiberService) {
	fiberSvcs = append(fiberSvcs, svc...)
}

var faberIsRegistered = false

func faberRegistered() {
	faberIsRegistered = true
	fiberSvcs = nil
}

func FiberApi(f func()) {
	if !faberIsRegistered {
		f()
	}
}

func fiberResHandler(ctx *fiber.Ctx, result []reflect.Value) error {
	writer := ctx.Response().BodyWriter()
	if !result[1].IsNil() {
		return json.NewEncoder(writer).Encode(errorcode.ErrHandle(result[1].Interface()))
	}
	if info, ok := result[0].Interface().(*http_fs.File); ok {
		header := ctx.Response().Header
		header.Set(httpi.HeaderContentType, httpi.ContentBinaryHeaderValue)
		header.Set(httpi.HeaderContentDisposition, "attachment;filename="+info.Name)
		io.Copy(writer, info.File)
		if flusher, canFlush := writer.(http.Flusher); canFlush {
			flusher.Flush()
		}
		return info.File.Close()
	}
	return ctx.JSON(httpi.ResAnyData{
		Code:    0,
		Message: "success",
		Details: result[0].Interface(),
	})
}

func FiberWithCtx(engine *fiber.App, genApi bool, modName string) {

	for _, v := range fiberSvcs {
		describe, preUrl, middleware := v.FiberService()
		value := reflect.ValueOf(v)
		if value.Kind() != reflect.Ptr {
			log.Fatal("必须传入指针")
		}
		var infos []*apiDocInfo
		engine.Group(preUrl, middleware...)
		for j := 0; j < value.NumMethod(); j++ {
			method := value.Type().Method(j)
			methodInfo := getMethodInfo(&method, preUrl, claimsType)
			if methodInfo == nil {
				continue
			}
			if methodInfo.path == "" || methodInfo.method == "" || methodInfo.title == "" || methodInfo.createlog.version == "" {
				log.Fatal("接口路径,方法,描述,创建日志均为必填")
			}
			methodType := method.Type
			methodValue := method.Func
			in2Type := methodType.In(2)
			engine.Add(methodInfo.method, methodInfo.path, func(ctx *fiber.Ctx) error {
				in1 := reflect.ValueOf(fasthttp_context.CtxWithRequest(context.Background(), ctx.Request()))
				in2 := reflect.New(in2Type.Elem())
				if err := fiber_build.Bind(ctx, in2.Interface()); err != nil {
					return ctx.Status(http.StatusBadRequest).JSON(errorcode.InvalidArgument.ErrRep())
				}
				result := methodValue.Call([]reflect.Value{value, in1, in2})
				return fiberResHandler(ctx, result)
			})
			infos = append(infos, &apiDocInfo{methodInfo, method.Type})
		}
		groupApiInfos = append(groupApiInfos, &groupApiInfo{describe, infos})
	}
	if genApi {
		filePath := apidoc.FilePath
		md(filePath, modName)
		swagger(filePath, modName)
		//gin_build.OpenApi(engine, filePath)
	}
	faberRegistered()
}
