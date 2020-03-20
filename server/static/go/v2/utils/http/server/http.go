package server

import (
	"net/http"

	"github.com/kataras/iris/v12"
	"github.com/liov/hoper/go/v2/initialize"
	iris_build "github.com/liov/hoper/go/v2/utils/http/iris"
	"github.com/liov/hoper/go/v2/utils/http/iris/api"
	"github.com/liov/hoper/go/v2/utils/http/iris/gateway"
	iris_log "github.com/liov/hoper/go/v2/utils/http/iris/log"
	"github.com/liov/hoper/go/v2/utils/log"
)

func (s *Server) Http() http.Handler {
	irisHandle := func(mux *iris.Application) {
		iris_build.WithConfiguration(mux, initialize.ConfUrl)
		logger := (&log.Config{Development: s.Conf.GetBasicConfig().Env == initialize.PRODUCT}).NewLogger()
		iris_log.SetLog(mux, logger, false)
		api.OpenApi(mux, "../protobuf/api/")
	}
	mux := iris_gateway.Http(irisHandle, s.HTTPRegistr)
	return mux
}
