// Copyright 2013 Julien Schmidt. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

// Package httprouter is a trie based high performance HTTP request router.
//
// A trivial example is:
//
//	package main
//
//	import (
//	    "fmt"
//	    "github.com/julienschmidt/httprouter"
//	    "net/http"
//	    "log"
//	)
//
//	func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
//	    fmt.Fprint(w, "Welcome!\n")
//	}
//
//	func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	    fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
//	}
//
//	func main() {
//	    router := httprouter.New()
//	    router.GET("/", Index)
//	    router.GET("/hello/:name", Hello)
//
//	    log.Fatal(http.ListenAndServe(":8080", router))
//	}
//
// The router matches incoming requests by the request method and the path.
// If a handle is registered for this path and method, the router delegates the
// request to that function.
// For the methods GET, POST, PUT, PATCH, DELETE and OPTIONS shortcut functions exist to
// register handles, for all other methods router.Handle can be used.
//
// The registered path, against which the router matches incoming requests, can
// contain two types of parameters:
//
//	Syntax    Type
//	:name     named parameter
//	*name     catch-all parameter
//
// Named parameters are dynamic path segments. They match anything until the
// next '/' or the path end:
//
//	Path: /blog/:category/:post
//
//	Requests:
//	 /blog/go/request-routers            match: category="go", post="request-routers"
//	 /blog/go/request-routers/           no match, but the router would redirect
//	 /blog/go/                           no match
//	 /blog/go/request-routers/comments   no match
//
// Catch-all parameters match anything until the path end, including the
// directory index (the '/' before the catch-all). Since they match anything
// until the end, catch-all parameters must always be the final path element.
//
//	Path: /files/*filepath
//
//	Requests:
//	 /files/                             match: filepath="/"
//	 /files/LICENSE                      match: filepath="/LICENSE"
//	 /files/templates/article.html       match: filepath="/templates/article.html"
//	 /files                              no match, but the router would redirect
//
// The value of parameters is saved as a slice of the Param struct, consisting
// each of a key and a value. The slice is passed to the Handle func as a third
// parameter.
// There are two ways to retrieve the value of a parameter:
//
//	// by the name of the parameter
//	user := ps.ByName("user") // defined by :user or *user
//
//	// by the index of the parameter. This way you can also get the name (key)
//	thirdKey   := ps[2].Key   // the name of the 3rd parameter
//	thirdValue := ps[2].Value // the value of the 3rd parameter
package pick

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/liov/hoper/server/go/lib/utils/concolor"
	httpi "github.com/liov/hoper/server/go/lib/utils/net/http"
	"github.com/liov/hoper/server/go/lib/utils/strings"
)

// Param is a single URL parameter, consisting of a key and a value.
type Param struct {
	Key   string
	Value string
}

// Params is a Param-slice, as returned by the router.
// The slice is ordered, the first URL parameter is also the first slice value.
// It is therefore safe to read values by the index.
type Params []Param

// ByName returns the value of the first Param which key matches the given name.
// If no matching Param is found, an empty string is returned.
func (ps Params) ByName(name string) string {
	for _, p := range ps {
		if p.Key == name {
			return p.Value
		}
	}
	return ""
}

type paramsKey struct{}

// ParamsKey is the request context key under which URL params are stored.
var ParamsKey = paramsKey{}

// ParamsFromContext pulls the URL parameters from a request context,
// or returns nil if none are present.
func ParamsFromContext(ctx context.Context) Params {
	p, _ := ctx.Value(ParamsKey).(Params)
	return p
}

// MatchedRoutePathParam is the Param name under which the path of the matched
// route is stored, if Router.SaveMatchedRoutePath is set.
var MatchedRoutePathParam = "$matchedRoutePath"

// MatchedRoutePath retrieves the path of the matched route.
// Router.SaveMatchedRoutePath must have been enabled when the respective
// handler was added, otherwise this function always returns an empty string.
func (ps Params) MatchedRoutePath() string {
	return ps.ByName(MatchedRoutePathParam)
}

// Router is a http.Handler which can be used to dispatch requests to different
// handler functions via configurable routes
type Router struct {
	trees *node

	paramsPool sync.Pool
	maxParams  uint16

	//前后调用
	middleware httpi.HandlerFuncs

	// If enabled, adds the matched route path onto the http.Request context
	// before invoking the handler.
	// The matched route path is only added to handlers of routes that were
	// registered when this option was enabled.
	SaveMatchedRoutePath bool

	// Enables automatic redirection if the current route can't be matched but a
	// handler for the path with (without) the trailing slash exists.
	// For example if /foo/ is requested but a route only exists for /foo, the
	// client is redirected to /foo with http status code 301 for GET requests
	// and 308 for all other request methods.
	RedirectTrailingSlash bool

	// If enabled, the router tries to fix the current request path, if no
	// handle is registered for it.
	// First superfluous path elements like ../ or // are removed.
	// Afterwards the router does a case-insensitive lookup of the cleaned path.
	// If a handle can be found for this route, the router makes a redirection
	// to the corrected path with status code 301 for GET requests and 308 for
	// all other request methods.
	// For example /FOO and /..//Foo could be redirected to /foo.
	// RedirectTrailingSlash is independent of this option.
	RedirectFixedPath bool

	// If enabled, the router checks if another method is allowed for the
	// current route, if the current request can not be routed.
	// If this is the case, the request is answered with 'Method Not Allowed'
	// and HTTP status code 405.
	// If no other Method is allowed, the request is delegated to the NotFound
	// handler.
	HandleMethodNotAllowed bool

	// If enabled, the router automatically replies to OPTIONS requests.
	// Customize OPTIONS handlers take priority over automatic replies.
	HandleOPTIONS bool

	// An optional http.Handler that is called on automatic OPTIONS requests.
	// The handler is only called if HandleOPTIONS is true and no OPTIONS
	// handler for the specific path was set.
	// The "Allowed" header is set before calling the handler.
	GlobalOPTIONS http.Handler

	// Cached value of global (*) allowed methods
	globalAllowed string

	// Configurable http.Handler which is called when no matching route is
	// found. If it is not set, http.NotFound is used.
	NotFound http.Handler

	// Configurable http.Handler which is called when a request
	// cannot be routed and HandleMethodNotAllowed is true.
	// If it is not set, http.Error with http.StatusMethodNotAllowed is used.
	// The "Allow" header with allowed request methods is set before the handler
	// is called.
	MethodNotAllowed http.Handler

	// Function to handle panics recovered from http handlers.
	// It should be used to generate a error page and return the http error code
	// 500 (Internal Server Error).
	// The handler can be used to keep your server from crashing because of
	// unrecovered panics.
	PanicHandler func(http.ResponseWriter, *http.Request, interface{})

	Tracing bool
}

// Make sure the Router conforms with the http.Handler interface
var _ http.Handler = func() *Router { return &Router{} }()

func (r *Router) getParams() *Params {
	ps := r.paramsPool.Get().(*Params)
	*ps = (*ps)[0:0] // reset slice
	return ps
}

func (r *Router) putParams(ps *Params) {
	if ps != nil {
		r.paramsPool.Put(ps)
	}
}

func New(genApi bool, modName string) *Router {
	router := &Router{
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      true,
		HandleMethodNotAllowed: true,
		HandleOPTIONS:          true,
		middleware:             make([]http.HandlerFunc, 0),
	}
	register(router, genApi, modName)
	return router
}

// Handle registers a new request handle with the given path and method.
//
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (r *Router) Handle(method, path string, middleware []http.HandlerFunc, handle reflect.Value) {
	varsCount := uint16(0)

	if method == "" {
		panic("method must not be empty")
	}
	if len(path) < 1 || path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}

	if handle.IsNil() {
		panic("handle must not be empty")
	}

	if r.SaveMatchedRoutePath {
		varsCount++
	}

	if r.trees == nil {
		r.trees = new(node)
	}

	r.trees.addRoute(path, &methodHandle{method, middleware, nil, &handle})

	// Update maxParams
	if paramsCount := countParams(path); paramsCount+varsCount > r.maxParams {
		r.maxParams = paramsCount + varsCount
	}

	// Lazy-init paramsPool alloc func
	if r.paramsPool.New == nil && r.maxParams > 0 {
		r.paramsPool.New = func() interface{} {
			ps := make(Params, 0, r.maxParams)
			return &ps
		}
	}
}

func (r *Router) Handler(method, path string, handle ...http.HandlerFunc) {
	varsCount := uint16(0)

	if method == "" {
		panic("method must not be empty")
	}
	if len(path) < 1 || path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}

	if handle == nil {
		panic("handle must not be empty")
	}

	if r.SaveMatchedRoutePath {
		varsCount++
	}

	if r.trees == nil {
		r.trees = new(node)
	}

	r.trees.addRoute(path, &methodHandle{method, handle[:len(handle)-1], handle[len(handle)-1], &reflect.Value{}})

	// Update maxParams
	if paramsCount := countParams(path); paramsCount+varsCount > r.maxParams {
		r.maxParams = paramsCount + varsCount
	}

	// Lazy-init paramsPool alloc func
	if r.paramsPool.New == nil && r.maxParams > 0 {
		r.paramsPool.New = func() interface{} {
			ps := make(Params, 0, r.maxParams)
			return &ps
		}
	}
}

// ServeFiles serves files from the given file system root.
// The path must end with "/*filepath", files are then served from the local
// path /defined/root/dir/*filepath.
// For example if root is "/etc" and *filepath is "passwd", the local file
// "/etc/passwd" would be served.
// Internally a http.FileServer is used, therefore http.NotFound is used instead
// of the Router's NotFound handler.
// To use the operating system's file system implementation,
// use http.Dir:
//
//	router.ServeFiles("/src/*filepath", http.Dir("/var/www"))
func (r *Router) ServeFiles(path string, root string) {

	fileServer := http.FileServer(http.Dir(root))

	r.Handler(http.MethodGet, path+"/*filepath", nil, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		req.URL.Path = req.URL.Path[len(path):]
		fileServer.ServeHTTP(w, req)
	}))
}

func (r *Router) Use(middleware ...http.HandlerFunc) {
	r.middleware = append(r.middleware, middleware...)
}

func (r *Router) GroupUse(path string, middleware ...http.HandlerFunc) {
	if r.trees == nil {
		r.trees = new(node)
	}
	r.trees.use(path, middleware...)
}

func (r *Router) recv(w http.ResponseWriter, req *http.Request) {
	if rcv := recover(); rcv != nil {
		r.PanicHandler(w, req, rcv)
	}
}

func (r *Router) allowed(path, reqMethod string, mhs []*methodHandle) (allow string) {
	allowed := make([]string, 0, 9)

	if path == "*" { // server-wide
		// empty method is used for internal calls to refresh the cache
		return r.globalAllowed
	} else { // specific path
		for _, mh := range mhs {
			if mh.method == http.MethodOptions {
				continue
			}
			// Add request method to list of allowed methods
			allowed = append(allowed, mh.method)
		}

	}

	return allowedMethod(allowed)
}

func allowedMethod(allowed []string) string {
	if len(allowed) > 0 {
		// Add request method to list of allowed methods
		allowed = append(allowed, http.MethodOptions)

		// Sort allowed methods.
		// sort.Strings(allowed) unfortunately causes unnecessary allocations
		// due to allowed being moved to the heap and interface conversion
		for i, l := 1, len(allowed); i < l; i++ {
			for j := i; j > 0 && allowed[j] < allowed[j-1]; j-- {
				allowed[j], allowed[j-1] = allowed[j-1], allowed[j]
			}
		}

		// return as comma separated list
		return strings.Join(allowed, ", ")
	}
	return ""
}

// ServeHTTP makes the router implement the http.Handler interface.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if r.PanicHandler != nil {
		defer r.recv(w, req)
	}

	path := req.URL.Path

	if root := r.trees; root != nil {
		r.middleware.ServeHTTP(w, req)
		if middleware, handles, ps, tsr := root.getValue(path, r.getParams); handles != nil {
			mh := getHandle(req.Method, handles)
			if mh.Valid() {
				for i := range middleware {
					middleware[i](w, req)
				}
				for i := range mh.middleware {
					mh.middleware[i](w, req)
				}
			}
			if mh.httpHandler != nil {
				mh.httpHandler.ServeHTTP(w, req)
				return
			}
			if mh.handle.IsValid() {
				if r.SaveMatchedRoutePath {
					*ps = append(*ps, Param{Key: MatchedRoutePathParam, Value: path})
				}
				commonHandler(w, req, mh.handle, ps, r.Tracing)
				return
			}
			if allow := r.allowed(path, req.Method, handles); allow != "" {
				w.Header().Set("Allow", allow)
				if req.Method == http.MethodOptions && r.HandleOPTIONS && r.GlobalOPTIONS != nil {
					r.GlobalOPTIONS.ServeHTTP(w, req)
				} else if r.HandleMethodNotAllowed {
					if r.MethodNotAllowed != nil {
						r.MethodNotAllowed.ServeHTTP(w, req)
					} else {
						http.Error(w,
							http.StatusText(http.StatusMethodNotAllowed),
							http.StatusMethodNotAllowed,
						)
					}
				}
				return
			}
		} else if req.Method != http.MethodConnect && path != "/" {
			// Moved Permanently, request with GET method
			code := http.StatusMovedPermanently
			if req.Method != http.MethodGet {
				// Permanent Redirect, request with same method
				code = http.StatusPermanentRedirect
			}

			if r.RedirectTrailingSlash && tsr != nil && getHandle(req.Method, tsr.handle).Valid() {
				if len(path) > 1 && path[len(path)-1] == '/' {
					req.URL.Path = path[:len(path)-1]
				} else {
					req.URL.Path = path + "/"
				}
				http.Redirect(w, req, req.URL.String(), code)
				return
			}

			// Try to fix the request path
			if r.RedirectFixedPath {
				fixedPath, found := root.findCaseInsensitivePath(
					httpi.CleanPath(path),
					req.Method,
					r.RedirectTrailingSlash,
				)
				if found {
					req.URL.Path = fixedPath
					http.Redirect(w, req, req.URL.String(), code)
					return
				}
			}
		}
	}

	// Handle 404
	if r.NotFound != nil {
		r.NotFound.ServeHTTP(w, req)
	} else {
		http.NotFound(w, req)
	}
}

func Log(method, path, title string) {
	fmt.Printf(" %s\t %s %s\t %s\n",
		concolor.Green("API:"),
		concolor.Yellow(stringsi.FormatLen(method, 6)),
		concolor.Blue(stringsi.FormatLen(path, 50)), concolor.Purple(title))
}
