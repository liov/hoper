package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/actliboy/hoper/server/go/lib/utils/log"
	neti "github.com/actliboy/hoper/server/go/lib/utils/net"
	mhttp "github.com/actliboy/hoper/server/go/lib/utils/net/http"
	"github.com/actliboy/hoper/server/go/lib/utils/number"
	"github.com/actliboy/hoper/server/go/lib/utils/strings"
	"go.uber.org/zap"
)

var client = &http.Client{}

const timeout = time.Minute

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

func SetTimeout(timeout time.Duration) {
	if timeout < time.Second {
		timeout = timeout * time.Second
	}
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

type LogCallback func(url, method, auth, reqBody, respBytes string, status int, process time.Duration)

func defaultLog(url, method, auth, reqBody, respBytes string, status int, process time.Duration) {
	log.Default.Logger.Info("third-request", zap.String("interface", url),
		zap.String("method", method),
		zap.String("param", reqBody),
		zap.Duration("processTime", process),
		zap.String("result", respBytes),
		zap.String("other", auth),
		zap.Int("status", status),
		zap.String("source", neti.GetIP()))
}

type ContentType uint8

const (
	ContentTypeJson     ContentType = iota
	ContentTypeForm     ContentType = iota
	ContentTypeFormData ContentType = iota
)

// RequestParams ...
type RequestParams struct {
	url, method        string
	timeout            time.Duration
	AuthUser, AuthPass string
	ContentType        ContentType
	Param              interface{}
	Header             http.Header
	logger             LogCallback
}

func NewRequest(url, method string, param interface{}) *RequestParams {
	return &RequestParams{url: url, method: strings.ToUpper(method), Header: make(http.Header), Param: param, logger: defaultLog}
}

func (req *RequestParams) SetParam(param interface{}) *RequestParams {
	req.Param = param
	return req
}

func (req *RequestParams) SetContentType(contentType ContentType) *RequestParams {
	req.ContentType = contentType
	return req
}

func (req *RequestParams) SetHeader(k, v string) *RequestParams {
	req.Header.Set(k, v)
	return req
}

func (req *RequestParams) SetLogger(logger LogCallback) *RequestParams {
	req.logger = logger
	return req
}

func (req *RequestParams) SetTimeout(timeout time.Duration) *RequestParams {
	req.timeout = timeout
	return req
}

type responseBody interface {
	CheckError() error
}

type ResponseBody struct {
	Status  int         `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func CommonResponse(response interface{}) responseBody {
	return &ResponseBody{Data: response}
}

func (res *ResponseBody) CheckError() error {
	if res.Status != 0 {
		return errors.New(res.Message)
	}
	return nil
}

// Do create a HTTP request
func (req *RequestParams) Do(response interface{}) error {
	method := req.method
	url := req.url
	if req.timeout != 0 {
		defer SetTimeout(timeout)
		SetTimeout(req.timeout)
	}
	var body io.Reader
	var reqBody, respBody string
	var statusCode int
	// 日志记录
	defer func(now time.Time) {
		if req.logger != nil {
			req.logger(url, method, req.AuthUser, reqBody, respBody, statusCode, time.Now().Sub(now))
		}
	}(time.Now())

	var err error
	if method == http.MethodGet {
		if req.Param != nil {
			param := getParam(req.Param)
			reqBody = param
			url += "?" + param
		}
	} else {
		if req.ContentType == ContentTypeJson {
			reqBytes, err := json.Marshal(req.Param)
			if err != nil {
				return err
			}
			body = bytes.NewReader(reqBytes)
			reqBody = stringsi.ToString(reqBytes)
		} else {
			param := getParam(req.Param)
			reqBody = param
			body = strings.NewReader(param)
		}
	}

	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	request.Header = req.Header
	if req.AuthUser != "" && req.AuthPass != "" {
		request.SetBasicAuth(req.AuthUser, req.AuthPass)
	}
	if req.ContentType == ContentTypeJson {
		request.Header.Set("Content-Type", "application/json;charset=utf-8")
	} else if req.ContentType == ContentTypeFormData {
		request.Header.Set("Content-Type", mhttp.ContentFormMultipartHeaderValue)
	} else {
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded;param=value")
	}

	resp, err := client.Do(request)
	if err != nil {
		respBody = err.Error()
		return err
	}
	defer resp.Body.Close()

	if resp.ContentLength != 0 && response != nil {
		respBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		respBody = stringsi.ToString(respBytes)
		statusCode = resp.StatusCode
		if resp.StatusCode < 200 || resp.StatusCode > 300 {
			return errors.New("status:" + resp.Status + respBody)
		}
		if len(respBytes) > 0 {
			err = json.Unmarshal(respBytes, response)
			if err != nil {
				return err
			}
			if v, ok := response.(responseBody); ok {
				err = v.CheckError()
			}
		}
	}

	return err
}

func getParam(param interface{}) string {
	if param == nil {
		return ""
	}
	query := url.Values{}
	parseParam(param, query)
	return query.Encode()
}

func parseParam(param interface{}, query url.Values) {
	v := reflect.ValueOf(param).Elem()
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		filed := v.Field(i)
		kind := filed.Kind()
		if kind == reflect.Interface || kind == reflect.Ptr {
			parseParam(filed.Interface(), query)
			continue
		}
		if kind == reflect.Struct {
			parseParam(filed.Addr().Interface(), query)
			continue
		}
		value := getFieldValue(filed)
		if value != "" {
			query.Set(t.Field(i).Tag.Get("json"), getFieldValue(v.Field(i)))
		}
	}

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
	case reflect.Struct:

	}
	return ""
}

type EasyRequest = RequestParams

func NewEasyRequest() *EasyRequest {
	return &EasyRequest{Header: make(http.Header), logger: defaultLog}
}

func Get(url string, param interface{}) *RequestParams {
	return NewRequest(url, http.MethodGet, param)
}

func (req *EasyRequest) Get(url string) *EasyRequest {
	req.url = url
	req.method = http.MethodGet
	return req
}

func (req *EasyRequest) DoGet(url string, param, response interface{}) error {
	req.url = url
	req.Param = param
	req.method = http.MethodGet
	return (*RequestParams)(req).Do(response)
}

func Post(url string, param interface{}) *RequestParams {
	return NewRequest(url, http.MethodPost, param)
}

func (req *RequestParams) Post(url string) *RequestParams {
	req.url = url
	req.method = http.MethodPost
	return req
}

func (req *EasyRequest) DoPost(url string, param, response interface{}) error {
	req.url = url
	req.Param = param
	req.method = http.MethodPost
	return (*RequestParams)(req).Do(response)
}

func Put(url string, param interface{}) *RequestParams {
	return NewRequest(url, http.MethodPut, param)
}

func (req *RequestParams) Put(url string) *RequestParams {
	req.url = url
	req.method = http.MethodPut
	return req
}

func (req *EasyRequest) DoPut(url string, param, response interface{}) error {
	req.url = url
	req.Param = param
	req.method = http.MethodPut
	return (*RequestParams)(req).Do(response)
}

func Delete(url string, param interface{}) *RequestParams {
	return NewRequest(url, http.MethodDelete, param)
}

func (req *RequestParams) Delete(url string) *RequestParams {
	req.url = url
	req.method = http.MethodDelete
	return req
}

func (req *EasyRequest) DoDelete(url string, param, response interface{}) error {
	req.url = url
	req.Param = param
	req.method = http.MethodDelete
	return (*RequestParams)(req).Do(response)
}

func (req *EasyRequest) CompleteDo(url, method string, param, response interface{}) error {
	req.url = url
	req.Param = param
	req.method = method
	return (*RequestParams)(req).Do(response)
}
