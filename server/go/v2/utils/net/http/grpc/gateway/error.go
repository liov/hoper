package gateway

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/liov/hoper/go/v2/protobuf/utils/errorcode"
	"github.com/liov/hoper/go/v2/utils/net/http/grpc/reconn"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

func CustomHTTPError(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {

	s, ok := status.FromError(err)
	if ok && s.Code() == 14 && strings.HasSuffix(s.Message(), `refused it."`) {
		//提供一个思路，这里应该是哪条连接失败重连哪条，不能这么粗暴，map的key是个关键
		if len(reconn.ReConnectMap) > 0 {
			for _, f := range reconn.ReConnectMap {
				f()
			}
		}
	}

	const fallback = `{"code": 14, "message": "failed to marshal error message"}`

	w.Header().Del("Trailer")
	contentType := marshaler.ContentType()
	w.Header().Set("Content-Type", contentType)
	se, ok := err.(*errorcode.ErrRep)
	if !ok {
		se = &errorcode.ErrRep{Code: errorcode.Unknown, Message: err.Error()}
	}

	buf, merr := marshaler.Marshal(se)
	if merr != nil {
		grpclog.Infof("Failed to marshal error message %q: %v", se, merr)
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := io.WriteString(w, fallback); err != nil {
			grpclog.Infof("Failed to write response: %v", err)
		}
		return
	}

	st := HTTPStatusFromCode(se.Code)
	w.WriteHeader(st)
	if _, err := w.Write(buf); err != nil {
		grpclog.Infof("Failed to write response: %v", err)
	}
}

func HTTPStatusFromCode(code errorcode.ErrCode) int {
	switch code {
	case errorcode.SUCCESS:
		return http.StatusOK
	case errorcode.Canceled:
		return http.StatusRequestTimeout
	case errorcode.Unknown:
		return http.StatusInternalServerError
	case errorcode.InvalidArgument:
		return http.StatusBadRequest
	case errorcode.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case errorcode.NotFound:
		return http.StatusNotFound
	case errorcode.AlreadyExists:
		return http.StatusConflict
	case errorcode.PermissionDenied:
		return http.StatusForbidden
	case errorcode.Unauthenticated:
		return http.StatusUnauthorized
	case errorcode.ResourceExhausted:
		return http.StatusTooManyRequests
	case errorcode.FailedPrecondition:
		// Note, this deliberately doesn't translate to the similarly named '412 Precondition Failed' HTTP response status.
		return http.StatusBadRequest
	case errorcode.Aborted:
		return http.StatusConflict
	case errorcode.OutOfRange:
		return http.StatusBadRequest
	case errorcode.Unimplemented:
		return http.StatusNotImplemented
	case errorcode.Internal:
		return http.StatusInternalServerError
	case errorcode.Unavailable:
		return http.StatusServiceUnavailable
	case errorcode.DataLoss:
		return http.StatusInternalServerError
	}

	grpclog.Infof("Unknown gRPC error code: %v", code)
	return http.StatusInternalServerError
}
