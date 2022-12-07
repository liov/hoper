package httpi

const (
	// ContentBinaryHeaderValue header value for binary data.
	ContentBinaryHeaderValue = "application/octet-stream"
	// ContentWebassemblyHeaderValue header value for web assembly files.
	ContentWebassemblyHeaderValue = "application/wasm"
	// ContentHTMLHeaderValue is the  string of text/html response header's content type value.
	ContentHTMLHeaderValue = "text/html"
	// ContentJSONHeaderValue header value for JSON data.
	ContentJSONHeaderValue = "application/json"
	// ContentJSONProblemHeaderValue header value for JSON API problem error.
	// Read more at: https://tools.ietf.org/html/rfc7807
	ContentJSONProblemHeaderValue = "application/problem+json"
	// ContentXMLProblemHeaderValue header value for XML API problem error.
	// Read more at: https://tools.ietf.org/html/rfc7807
	ContentXMLProblemHeaderValue = "application/problem+xml"
	// ContentJavascriptHeaderValue header value for JSONP & Javascript data.
	ContentJavascriptHeaderValue = "text/javascript"
	// ContentTextHeaderValue header value for Text data.
	ContentTextHeaderValue = "text/plain"
	// ContentXMLHeaderValue header value for XML data.
	ContentXMLHeaderValue = "text/xml"
	// ContentXMLUnreadableHeaderValue obselete header value for XML.
	ContentXMLUnreadableHeaderValue = "application/xml"
	// ContentMarkdownHeaderValue custom key/content type, the real is the text/html.
	ContentMarkdownHeaderValue = "text/markdown"
	// ContentYAMLHeaderValue header value for YAML data.
	ContentYAMLHeaderValue = "application/x-yaml"
	// ContentYAMLTextHeaderValue header value for YAML plain text.
	ContentYAMLTextHeaderValue = "text/yaml"
	// ContentProtobufHeaderValue header value for Protobuf messages data.
	ContentProtobufHeaderValue = "application/x-protobuf"
	// ContentMsgPackHeaderValue header value for MsgPack data.
	ContentMsgPackHeaderValue = "application/msgpack"
	// ContentMsgPack2HeaderValue alternative header value for MsgPack data.
	ContentMsgPack2HeaderValue = "application/x-msgpack"
	// ContentFormHeaderValue header value for post form data.
	ContentFormHeaderValue = "application/x-www-form-urlencoded"
	// ContentFormMultipartHeaderValue header value for post multipart form data.
	ContentFormMultipartHeaderValue = "multipart/form-data"
	// ContentGRPCHeaderValue Content-Type header value for gRPC.
	ContentGRPCHeaderValue    = "application/grpc"
	ContentGRPCWebHeaderValue = "application/grpc-web"

	ContentJSONUTF8HeaderValue = "application/json;charset=utf-8"

	ContentFormParamHeaderValue = "application/x-www-form-urlencoded;param=value"
)

const (
	HeaderDeviceInfo                  = "Device-AuthInfo"
	HeaderLocation                    = "Location"
	HeaderArea                        = "Area"
	HeaderUserAgent                   = "User-Agent"
	HeaderXForwardedFor               = "X-Forwarded-For"
	HeaderAuth                        = "HeaderAuth"
	HeaderContentType                 = "Content-Type"
	HeaderTrace                       = "Trace"
	HeaderTraceID                     = "Trace-ID"
	HeaderTraceBin                    = "Trace-Bin"
	HeaderAuthorization               = "Authorization"
	HeaderCookie                      = "Cookie"
	HeaderCookieToken                 = "token"
	HeaderCookieDel                   = "del"
	HeaderContentDisposition          = "Content-Disposition"
	HeaderContentEncoding             = "Content-Encoding"
	HeaderReferer                     = "Referer"
	HeaderAccept                      = "Accept"
	HeaderCacheControl                = "Cache-Control"
	HeaderSetCookie                   = "Set-Cookie"
	HeaderTrailer                     = "Trailer"
	HeaderTransferEncoding            = "Transfer-Encoding"
	HeaderInternal                    = "Internal"
	HeaderTE                          = "TE"
	HeaderLastModified                = "Last-Modified"
	HeaderContentLength               = "Content-Length"
	HeaderAccessControlRequestMethod  = "Access-Control-Request-Method"
	HeaderAccessControlRequestHeaders = "Access-Control-Request-Headers"
)

const (
	GrpcTraceBin = "grpc-trace-bin"
	GrpcInternal = "grpc-internal"
)
