package stats

import (
	"context"
	httpi "github.com/liov/hoper/server/go/lib/utils/net/http"
	"go.opencensus.io/trace"
	"go.opencensus.io/trace/propagation"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
	"strings"
)

type ClientHandler struct {
	// StartOptions allows configuring the StartOptions used to create new spans.
	//
	// StartOptions.SpanKind will always be set to trace.SpanKindClient
	// for spans started by this handler.
	StartOptions trace.StartOptions
}

// HandleConn exists to satisfy gRPC stats.Handler.
func (c *ClientHandler) HandleConn(ctx context.Context, cs stats.ConnStats) {
	// no-op
}

// TagConn exists to satisfy gRPC stats.Handler.
func (c *ClientHandler) TagConn(ctx context.Context, cti *stats.ConnTagInfo) context.Context {
	// no-op
	return ctx
}

// HandleRPC implements per-RPC tracing and stats instrumentation.
func (c *ClientHandler) HandleRPC(ctx context.Context, rs stats.RPCStats) {
	span := trace.FromContext(ctx)
	// TODO: compressed and uncompressed sizes are not populated in every message.
	switch rs := rs.(type) {
	case *stats.Begin:
		span.AddAttributes(
			trace.BoolAttribute("Client", rs.Client),
			trace.BoolAttribute("FailFast", rs.FailFast))
	case *stats.InPayload:
		span.AddMessageReceiveEvent(0 /* TODO: messageID */, int64(rs.Length), int64(rs.WireLength))
	case *stats.OutPayload:
		span.AddMessageSendEvent(0, int64(rs.Length), int64(rs.WireLength))
	case *stats.End:
		if rs.Error != nil {
			s, ok := status.FromError(rs.Error)
			if ok {
				span.SetStatus(trace.Status{Code: int32(s.Code()), Message: s.Message()})
			} else {
				span.SetStatus(trace.Status{Code: int32(codes.Internal), Message: rs.Error.Error()})
			}
		}
		span.End()
	}
}

// TagRPC implements per-RPC context management.
func (c *ClientHandler) TagRPC(ctx context.Context, rti *stats.RPCTagInfo) context.Context {
	name := strings.TrimPrefix(rti.FullMethodName, "/")
	name = strings.Replace(name, "/", ".", -1)
	ctx, span := trace.StartSpan(ctx, name,
		trace.WithSampler(c.StartOptions.Sampler),
		trace.WithSpanKind(trace.SpanKindClient)) // span is ended by traceHandleRPC
	traceContextBinary := propagation.Binary(span.SpanContext())
	return metadata.AppendToOutgoingContext(ctx, httpi.GrpcTraceBin, string(traceContextBinary),
		httpi.GrpcInternal, httpi.GrpcInternal)
}
