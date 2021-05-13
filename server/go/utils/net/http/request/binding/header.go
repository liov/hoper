package binding

import (
	"net/http"
	"net/textproto"
	"reflect"

	"github.com/valyala/fasthttp"
)

type headerBinding struct{}

func (headerBinding) Name() string {
	return "header"
}

func (headerBinding) Bind(req *http.Request, obj interface{}) error {

	if err := mapHeader(obj, req.Header); err != nil {
		return err
	}

	return validate(obj)
}

func (headerBinding) FasthttpBind(req *fasthttp.Request, obj interface{}) error {

	if err := mappingByPtr(obj, (*reqSource)(&req.Header), tag); err != nil {
		return err
	}

	return validate(obj)
}

func mapHeader(ptr interface{}, h map[string][]string) error {
	return mappingByPtr(ptr, headerSource(h), tag)
}

type headerSource map[string][]string

var _ setter = headerSource(nil)

func (hs headerSource) Peek(key string) ([]string, bool) {
	v, ok := hs[key]
	return v, ok
}

func (hs headerSource) TrySet(value reflect.Value, field reflect.StructField, tagValue string, opt setOptions) (isSetted bool, err error) {
	return setByKV(value, field, hs, textproto.CanonicalMIMEHeaderKey(tagValue), opt)
}
