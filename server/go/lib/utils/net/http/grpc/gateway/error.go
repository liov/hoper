package gateway

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/textproto"
	"strings"

	"github.com/actliboy/hoper/server/go/lib/protobuf/errorcode"
	httpi "github.com/actliboy/hoper/server/go/lib/utils/net/http"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/grpc/reconn"
	stringsi "github.com/actliboy/hoper/server/go/lib/utils/strings"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func CustomHTTPError(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {

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

	w.Header().Del(httpi.HeaderTrailer)
	contentType := marshaler.ContentType(nil)
	w.Header().Set(httpi.HeaderContentType, contentType)
	se, ok := err.(*errorcode.ErrRep)
	if !ok {
		se = &errorcode.ErrRep{Code: errorcode.Unknown, Message: err.Error()}
	}

	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		grpclog.Infof("Failed to extract ServerMetadata from context")
	}

	handleForwardResponseServerMetadata(w, md.HeaderMD)

	buf, merr := marshaler.Marshal(se)
	if merr != nil {
		grpclog.Infof("Failed to marshal error message %q: %v", se, merr)
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := io.WriteString(w, fallback); err != nil {
			grpclog.Infof("Failed to write response: %v", err)
		}
		return
	}

	var wantsTrailers bool

	if te := r.Header.Get(httpi.HeaderTE); strings.Contains(strings.ToLower(te), "trailers") {
		wantsTrailers = true
		handleForwardResponseTrailerHeader(w, md.TrailerMD)
		w.Header().Set(httpi.HeaderTransferEncoding, "chunked")
	}

	/*	st := HTTPStatusFromCode(se.Code)
		w.WriteHeader(st)*/
	if _, err := w.Write(buf); err != nil {
		grpclog.Infof("Failed to write response: %v", err)
	}
	if wantsTrailers {
		handleForwardResponseTrailer(w, md.TrailerMD)
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

func outgoingHeaderMatcher(key string) (string, bool) {
	switch key {
	case
		httpi.HeaderSetCookie:
		return key, true
	}
	return "", false
}

func headerMatcher() []string {
	return []string{httpi.HeaderSetCookie}
}

func handleForwardResponseServerMetadata(w http.ResponseWriter, md metadata.MD) {
	for _, k := range headerMatcher() {
		if vs, ok := md[k]; ok {
			for _, v := range vs {
				w.Header().Add(k, v)
			}
		}
	}
}

func handleForwardResponseTrailerHeader(w http.ResponseWriter, md metadata.MD) {
	for k := range md {
		tKey := textproto.CanonicalMIMEHeaderKey(fmt.Sprintf("%s%s", runtime.MetadataTrailerPrefix, k))
		w.Header().Add("Trailer", tKey)
	}
}

func handleForwardResponseTrailer(w http.ResponseWriter, md metadata.MD) {
	for k, vs := range md {
		tKey := fmt.Sprintf("%s%s", runtime.MetadataTrailerPrefix, k)
		for _, v := range vs {
			w.Header().Add(tKey, v)
		}
	}
}

func RoutingErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, httpStatus int) {
	w.WriteHeader(httpStatus)
	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.Write(stringsi.ToBytes(http.StatusText(httpStatus)))
}
