package gini

import "github.com/gin-gonic/gin"

type Config struct {
	RedirectTrailingSlash bool

	RedirectFixedPath bool

	HandleMethodNotAllowed bool
	ForwardedByClientIP    bool

	AppEngine bool

	UseRawPath bool

	UnescapePathValues bool

	MaxMultipartMemory int64

	RemoveExtraSlash bool
}

func (c *Config) SetConfig(engine *gin.Engine) {
	if c == nil {
		return
	}
	engine.RedirectTrailingSlash = c.RedirectTrailingSlash
	engine.RedirectFixedPath = c.RedirectFixedPath
	engine.HandleMethodNotAllowed = c.HandleMethodNotAllowed
	engine.ForwardedByClientIP = c.ForwardedByClientIP
	engine.AppEngine = c.AppEngine
	engine.UseRawPath = c.UseRawPath
	engine.UnescapePathValues = c.UnescapePathValues
	engine.MaxMultipartMemory = c.MaxMultipartMemory
	engine.RemoveExtraSlash = c.RemoveExtraSlash
}
