package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/liov/hoper/go/v2/utils/number"
)

var client = http.DefaultClient

const timeout = 3 * time.Second

func init() {
	client.Transport = &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(network, addr, timeout)
			if err != nil {
				return nil, err
			}
			c.SetDeadline(time.Now().Add(timeout))
			return c, nil
		},
		DisableKeepAlives: true,
	}
}

type Pair struct {
	K, V string
}

type LogCallback func(string, string, []byte,int, []byte)

// RequestParams ...
type RequestParams struct {
	url, method        string
	Timeout            time.Duration
	AuthUser, AuthPass string
	ContentType        string
	Param              interface{}
	Header             []Pair
	logger             LogCallback
}

func NewRequest(url, method string, param interface{}) *RequestParams {
	return &RequestParams{url: url, method: strings.ToUpper(method), Param: param}
}

func (req *RequestParams) SetHeader(k, v string) *RequestParams {
	req.Header = append(req.Header, Pair{K: k, V: v})
	return req
}

func (req *RequestParams) SetLogger(logger LogCallback) *RequestParams {
	req.logger = logger
	return req
}

type responseBody interface {
	CheckError() error
}

type ResponseBody struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	Msg    string      `json:"msg"`
}

func CommonResponse(response interface{}) interface{} {
	return &ResponseBody{Data: response}
}

func (res *ResponseBody) CheckError() error {
	if res.Status != 0 {
		return errors.New(res.Msg)
	}
	return nil
}

// HTTPRequest create a HTTP request
func (req *RequestParams) HTTPRequest(response interface{}) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("utils: HTTPRequest error: %v", err)
		}
	}()
	method := req.method
	url := req.url
	if req.Timeout != 0 {
		if req.Timeout < time.Second {
			req.Timeout = req.Timeout * time.Second
		}
		client.Transport = &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(network, addr, req.Timeout)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(time.Now().Add(req.Timeout))
				return c, nil
			},
			DisableKeepAlives: true,
		}
	}
	var body io.Reader
	var reqBody []byte
	if method == http.MethodGet {
		url += "?" + getParam(req.Param)
	} else {
		if req.ContentType != "" {
			body = strings.NewReader(getParam(req.Param))
		} else {
			reqBody, err = json.Marshal(req.Param)
			if err != nil {
				return
			}
			body = bytes.NewReader(reqBody)
		}

	}

	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return
	}
	request.SetBasicAuth(req.AuthUser, req.AuthPass)
	if req.ContentType != "" {
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded;param=value")
	} else {
		request.Header.Set("Content-Type", "application/json;charset=utf-8")
	}
	if len(req.Header) != 0 {
		for _, pair := range req.Header {
			request.Header.Set(pair.K, pair.V)
		}
	}

	resp, err := client.Do(request)
	if err != nil {
		return
	}
	defer func() {
		err = resp.Body.Close()
	}()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if req.logger != nil {
		req.logger(url, method, reqBody,resp.StatusCode, respBytes)
	}
	if resp.StatusCode != 200 {
		return errors.New("status:" + resp.Status)
	}
	err = json.Unmarshal(respBytes, response)
	if err != nil {
		return
	}
	if v, ok := response.(responseBody); ok {
		err = v.CheckError()
	}
	return
}

func getParam(param interface{}) string {
	query := url.Values{}
	v := reflect.ValueOf(param).Elem()
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		value := getFieldValue(v.Field(i))
		if value != "" {
			query.Set(t.Field(i).Tag.Get("json"), getFieldValue(v.Field(i)))
		}
	}

	return query.Encode()
}

func getFieldValue(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.Itoa(int(v.Int()))
	case reflect.Float32, reflect.Float64:
		return number.FormatFloat(v.Float())
	case reflect.String:
		return v.String()
	case reflect.Interface, reflect.Ptr:
		return getFieldValue(v.Elem())
	}
	return ""
}
