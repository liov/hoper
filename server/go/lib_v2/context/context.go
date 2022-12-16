package context

import (
	contexti2 "github.com/liov/hoper/server/go/lib/context"
	contexti "github.com/liov/hoper/server/go/lib_v2/utils/context"
	"net/http"
)

type Context struct {
	contexti.RequestContext[http.Request, contexti2.Authorization]
}
