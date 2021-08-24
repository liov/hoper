package grpci

import (
	"github.com/liov/hoper/server/go/lib/utils/log"
	"google.golang.org/grpc/grpclog"
)

func init() {
	grpclog.SetLoggerV2(log.Default)
}
