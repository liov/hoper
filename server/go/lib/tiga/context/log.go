package contexti

import (
	"github.com/liov/hoper/server/go/lib/protobuf/errorcode"
	"github.com/liov/hoper/server/go/lib/utils/log"
	"go.uber.org/zap"
)

func (c *Ctx) Error(args ...interface{}) {
	args = append(args, zap.String(log.TraceId, c.TraceID))
	log.Default.Error(args...)
}

func (c *Ctx) HandleError(err error) {
	if err != nil {
		log.Default.Error(err.Error(), zap.String(log.TraceId, c.TraceID))
	}
}

func (c *Ctx) ErrorLog(err, originErr error, funcName string) error {
	// caller 用原始logger skip刚好
	log.Default.Error(originErr.Error(), zap.String(log.TraceId, c.TraceID), zap.Int(log.Type, errorcode.Code(err)), zap.String(log.Position, funcName))
	return err
}
