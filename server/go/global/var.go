package global

import (
	"github.com/hopeio/initialize"
	"go.opentelemetry.io/otel"
)

const ScopeName = "github.com/liov/hoper/server/go/user"

var (
	Global         = initialize.NewGlobal[*config, *dao]()
	Dao    *dao    = Global.Dao
	Conf   *config = Global.Config

	Tracer = otel.Tracer(ScopeName)
	Meter  = otel.Meter(ScopeName)
)
