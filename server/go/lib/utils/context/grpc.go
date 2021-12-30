package contexti

import (
	httpi "github.com/liov/hoper/server/go/lib/utils/net/http"
	"google.golang.org/grpc/metadata"
)

func (c *RequestContext) SetHeader(md metadata.MD) error {
	for k, v := range md {
		c.Request.Header[k] = v
	}
	if c.ServerTransportStream != nil {
		err := c.ServerTransportStream.SetHeader(md)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *RequestContext) SendHeader(md metadata.MD) error {
	for k, v := range md {
		c.Request.Header[k] = v
	}
	if c.ServerTransportStream != nil {
		err := c.ServerTransportStream.SendHeader(md)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *RequestContext) WriteHeader(k, v string) error {
	c.Request.Header[k] = []string{v}

	if c.ServerTransportStream != nil {
		err := c.ServerTransportStream.SendHeader(metadata.MD{k: []string{v}})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *RequestContext) SetCookie(v string) error {
	c.Request.Header[httpi.HeaderSetCookie] = []string{v}

	if c.ServerTransportStream != nil {
		err := c.ServerTransportStream.SendHeader(metadata.MD{httpi.HeaderSetCookie: []string{v}})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *RequestContext) SetTrailer(md metadata.MD) error {
	for k, v := range md {
		c.Request.Header[k] = v
	}
	if c.ServerTransportStream != nil {
		err := c.ServerTransportStream.SetTrailer(md)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *RequestContext) Method() string {
	if c.ServerTransportStream != nil {
		return c.ServerTransportStream.Method()
	}
	return ""
}
