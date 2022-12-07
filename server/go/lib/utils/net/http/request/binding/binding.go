package binding

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
	stringsi "github.com/liov/hoper/server/go/lib/utils/strings"
	"github.com/liov/hoper/server/go/lib/utils/verification/validator"
	"github.com/valyala/fasthttp"
)

const (
	MIMEJSON              = "application/json"
	MIMEHTML              = "text/html"
	MIMEXML               = "application/xml"
	MIMEXML2              = "text/xml"
	MIMEPlain             = "text/plain"
	MIMEPOSTForm          = "application/x-www-form-urlencoded"
	MIMEMultipartPOSTForm = "multipart/form-data"
	MIMEPROTOBUF          = "application/x-protobuf"
	MIMEMSGPACK           = "application/x-msgpack"
	MIMEMSGPACK2          = "application/msgpack"
	MIMEYAML              = "application/x-yaml"
	tag                   = "json"
)

// Binding describes the interface which needs to be implemented for binding the
// data present in the request such as JSON request body, query parameters or
// the form POST.
type Binding interface {
	Name() string
	Bind(*http.Request, interface{}) error
	GinBind(*gin.Context, interface{}) error
	FasthttpBind(*fasthttp.Request, interface{}) error
	FiberBind(*fiber.Ctx, interface{}) error
}

// BindingBody adds BindBody method to Binding. BindBody is similar with GinBind,
// but it reads the body from supplied bytes instead of req.Body.
type BindingBody interface {
	Binding
	BindBody([]byte, interface{}) error
}

// BindingUri adds BindUri method to Binding. BindUri is similar with GinBind,
// but it read the Params.
type BindingUri interface {
	Name() string
	BindUri(map[string][]string, interface{}) error
}

// Validator is the default validator which implements the StructValidator
// interface. It uses https://github.com/go-playground/validator/tree/v8.18.2
// under the hood.
var Validator = validator.DefaultValidator

// These implement the Binding interface and can be used to bind the data
// present in the request to struct instances.
var (
	JSON          = jsonBinding{}
	XML           = xmlBinding{}
	Form          = formBinding{}
	Query         = queryBinding{}
	FormPost      = formPostBinding{}
	FormMultipart = formMultipartBinding{}
	ProtoBuf      = protobufBinding{}
	MsgPack       = msgpackBinding{}
	YAML          = yamlBinding{}
	Uri           = uriBinding{}
	Header        = headerBinding{}
)

// Default returns the appropriate Binding instance based on the HTTP method
// and the content type.
func Default(method string, contentType []byte) Binding {
	if method == http.MethodGet {
		return Query
	}

	switch stringsi.ToString(contentType) {
	case MIMEJSON:
		return JSON
	case MIMEXML, MIMEXML2:
		return XML
	case MIMEPROTOBUF:
		return ProtoBuf
	case MIMEMSGPACK, MIMEMSGPACK2:
		return MsgPack
	case MIMEYAML:
		return YAML
	case MIMEMultipartPOSTForm:
		return FormMultipart
	default: // case MIMEPOSTForm:
		return Form
	}
}

func validate(obj interface{}) error {
	return validator.Validate(obj)
}
