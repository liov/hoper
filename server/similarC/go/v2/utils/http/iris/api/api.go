package api

import (
	"mime"
	"net/http"

	"github.com/go-openapi/spec"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/core/handlerconv"
	"github.com/liov/hoper/go/v2/utils/http/api"
)

func OpenApi(mux *iris.Application) {
	_ = mime.AddExtensionType(".svg", "image/svg+xml")

	mux.Get(api.PrefixUri+"{mod:path}", handlerconv.FromStd(api.HttpHandle))
}

func ApiMiddle(ctx context.Context) {
	currentRouteName := ctx.GetCurrentRoute().Name()[len(ctx.Method()):]

	pathItem := new(spec.PathItem)

	api.GetDoc("../")
	doc := &api.Doc
	if doc.Paths != nil && doc.Paths.Paths != nil {
		if path, ok := doc.Paths.Paths[currentRouteName]; ok {
			pathItem = &path
		} else {
			pathItem = api.NewPathItem()
		}
	} else {
		doc.Paths = &spec.Paths{Paths: map[string]spec.PathItem{}}
		pathItem = api.NewPathItem()
	}

	parameters := make([]spec.Parameter, ctx.Params().Store.Len(), ctx.Params().Store.Len())

	params := ctx.Params().Store

	for i := range params {
		key := params[i].Key

		//val := params[i].ValueRaw
		parameters[i] = spec.Parameter{
			ParamProps: spec.ParamProps{
				Name:        key,
				In:          "path",
				Description: "Description",
			},
		}
	}

	if ctx.URLParam("apidoc") == "stop" {
		defer api.WriteToFile("../")
	}

	var res spec.Responses
	op := spec.Operation{
		OperationProps: spec.OperationProps{
			Description: "Description",
			Consumes:    []string{"application/x-www-form-urlencoded"},
			Tags:        []string{"Tags"},
			Summary:     "Summary",
			ID:          "currentRouteName" + ctx.Method(),
			Parameters:  parameters,
			Responses:   &res,
		},
	}

	switch ctx.Method() {
	case http.MethodGet:
		pathItem.Get = &op
	case http.MethodPost:
		pathItem.Post = &op
	case http.MethodPut:
		pathItem.Put = &op
	case http.MethodDelete:
		pathItem.Delete = &op
	case http.MethodOptions:
		pathItem.Options = &op
	case http.MethodPatch:
		pathItem.Patch = &op
	case http.MethodHead:
		pathItem.Head = &op
	}
	doc.Paths.Paths[currentRouteName] = *pathItem
	ctx.Next()
}
