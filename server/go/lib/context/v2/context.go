package v2

import (
	contexti2 "github.com/actliboy/hoper/server/go/lib/context"
	contexti "github.com/actliboy/hoper/server/go/lib/utils/generics/context"
	"net/http"
)

type Context struct {
	contexti.RequestContext[http.Request, contexti2.Authorization]
}
