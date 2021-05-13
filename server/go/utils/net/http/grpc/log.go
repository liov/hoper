package grpci

import (
	"github.com/liov/hoper/v2/utils/log"
	"google.golang.org/grpc/grpclog"
)

func init() {
	grpclog.SetLoggerV2(log.Default)
}
