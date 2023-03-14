package http3

import (
	"github.com/quic-go/quic-go/http3"
	"net/http"
)

func GetClient() *http.Client {
	return &http.Client{
		Transport: &http3.RoundTripper{},
	}
}
