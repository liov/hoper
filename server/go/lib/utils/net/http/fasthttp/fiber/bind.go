package fiber_build

import (
	"github.com/gofiber/fiber/v2"
	"github.com/liov/hoper/server/go/lib/utils/net/http/request/binding"
)

func Bind(c *fiber.Ctx, obj interface{}) error {
	b := binding.Default(c.Method(), c.Request().Header.ContentType())
	return MustBindWith(c, obj, b)
}

// BindJSON is a shortcut for c.MustBindWith(obj, binding.JSON).
func BindJSON(c *fiber.Ctx, obj interface{}) error {
	return MustBindWith(c, obj, binding.JSON)
}

// BindXML is a shortcut for c.MustBindWith(obj, binding.BindXML).
func BindXML(c *fiber.Ctx, obj interface{}) error {
	return MustBindWith(c, obj, binding.XML)
}

// BindQuery is a shortcut for c.MustBindWith(obj, binding.Query).
func BindQuery(c *fiber.Ctx, obj interface{}) error {
	return MustBindWith(c, obj, binding.Query)
}

// BindYAML is a shortcut for c.MustBindWith(obj, binding.YAML).
func BindYAML(c *fiber.Ctx, obj interface{}) error {
	return MustBindWith(c, obj, binding.YAML)
}

// BindUri binds the passed struct pointer using binding.Uri.
// It will abort the request with HTTP 400 if any error occurs.
func BindUri(c *fiber.Ctx, obj interface{}) error {
	return ShouldBindUri(c, obj)
}

// MustBindWith binds the passed struct pointer using the specified binding engine.
// It will abort the request with HTTP 400 if any error occurs.
// See the binding package.
func MustBindWith(c *fiber.Ctx, obj interface{}, b binding.Binding) error {
	return ShouldBindWith(c, obj, b)
}

// ShouldBind checks the Content-Type to select a binding engine automatically,
// Depending the "Content-Type" header different bindings are used:
//
//	"application/json" --> JSON binding
//	"application/xml"  --> XML binding
//
// otherwise --> returns an error
// It parses the request's body as JSON if Content-Type == "application/json" using JSON or XML as a JSON input.
// It decodes the json payload into the struct specified as a pointer.
// Like c.GinBind() but this method does not set the response status code to 400 and abort if the json is not valid.
func ShouldBind(c *fiber.Ctx, obj interface{}) error {
	b := binding.Default(c.Method(), c.Request().Header.ContentType())
	return ShouldBindWith(c, obj, b)
}

// ShouldBindJSON is a shortcut for c.ShouldBindWith(obj, binding.JSON).
func ShouldBindJSON(c *fiber.Ctx, obj interface{}) error {
	return ShouldBindWith(c, obj, binding.JSON)
}

// ShouldBindXML is a shortcut for c.ShouldBindWith(obj, binding.XML).
func ShouldBindXML(c *fiber.Ctx, obj interface{}) error {
	return ShouldBindWith(c, obj, binding.XML)
}

// ShouldBindQuery is a shortcut for c.ShouldBindWith(obj, binding.Query).
func ShouldBindQuery(c *fiber.Ctx, obj interface{}) error {
	return ShouldBindWith(c, obj, binding.Query)
}

// ShouldBindYAML is a shortcut for c.ShouldBindWith(obj, binding.YAML).
func ShouldBindYAML(c *fiber.Ctx, obj interface{}) error {
	return ShouldBindWith(c, obj, binding.YAML)
}

// ShouldBindUri binds the passed struct pointer using the specified binding engine.
func ShouldBindUri(c *fiber.Ctx, obj interface{}) error {
	return binding.Uri.BindUri(c, obj)
}

// ShouldBindWith binds the passed struct pointer using the specified binding engine.
// See the binding package.
func ShouldBindWith(c *fiber.Ctx, obj interface{}, b binding.Binding) error {
	return b.FiberBind(c, obj)
}
