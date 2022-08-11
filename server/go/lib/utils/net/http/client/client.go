package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/actliboy/hoper/server/go/lib/utils/log"
	httpi "github.com/actliboy/hoper/server/go/lib/utils/net/http"
	"github.com/actliboy/hoper/server/go/lib/utils/number"
	"github.com/actliboy/hoper/server/go/lib/utils/strings"
	"go.uber.org/zap"
	urlpkg "net/url"
)

// 不是并发安全的

var (
	defaultClient = &http.Client{}
	genlog        = true
	headerMap     = sync.Map{}
)

const timeout = time.Minute

func init() {
	defaultClient.Transport = &http.Transport{
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
	setTimeout(defaultClient, timeout)
}

func DisableLog() {
	genlog = false
}

func SetDefaultLogger(logger LogCallback) {
	defaultLog = logger
}

func SetClient(client *http.Client) {
	client = defaultClient
}

func setTimeout(client *http.Client, timeout time.Duration) {
	if client == nil {
		client = defaultClient
	}
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

var defaultLog = DefaultLogger

func DefaultLogger(url, method, auth string, reqBody, respBody *Body, status int, process time.Duration, err error) {
	reqField, respField := zap.Skip(), zap.Skip()
	if reqBody != nil {
		key := "param"
		if reqBody.IsJson() {
			reqField = zap.Reflect(key, log.BytesJson(reqBody.Data))
		} else if reqBody.IsProtobuf() {
			reqField = zap.Binary(key, reqBody.Data)
		} else {
			reqField = zap.String(key, stringsi.ToString(reqBody.Data))
		}
	}
	if respBody != nil {
		key := "result"
		if respBody.IsJson() {
			respField = zap.Reflect(key, log.BytesJson(respBody.Data))
		} else if respBody.IsProtobuf() {
			respField = zap.Binary(key, respBody.Data)
		} else {
			respField = zap.String(key, stringsi.ToString(respBody.Data))
		}
	}

	log.Default.Logger.Info("third-request", zap.String("interface", url),
		zap.String("method", method),
		reqField,
		zap.Duration("processTime", process),
		respField,
		zap.String("other", auth),
		zap.Int("status", status),
		zap.Error(err))
}

type ContentType uint8

const (
	ContentTypeJson     ContentType = iota
	ContentTypeForm     ContentType = iota
	ContentTypeFormData ContentType = iota
	ContentTypeProtobuf ContentType = iota
)

// RequestParams ...
type RequestParams struct {
	ctx                context.Context
	client             *http.Client
	url, method        string
	timeout            time.Duration
	AuthUser, AuthPass string
	ContentType        ContentType
	header             http.Header
	cachedHeaderKey    string
	logger             LogCallback
	genlog             bool
	ResponseHandler    func(response []byte) ([]byte, error)
	retryTimes         int
	retryHandle        func(*RequestParams)
}

func New(url string) *RequestParams {
	return newRequest(url, "")
}

func NewRequest(url, method string) *RequestParams {
	return newRequest(url, strings.ToUpper(method))
}

func newRequest(url, method string) *RequestParams {
	return &RequestParams{ctx: context.Background(), client: defaultClient, url: url, method: method, header: make(http.Header), logger: defaultLog}
}

func NewGetRequest(url string) *RequestParams {
	return newRequest(url, http.MethodGet)
}

func NewPostRequest(url string) *RequestParams {
	return newRequest(url, http.MethodPost)
}

func NewPutRequest(url string) *RequestParams {
	return newRequest(url, http.MethodPut)
}

func NewDeleteRequest(url string) *RequestParams {
	return newRequest(url, http.MethodDelete)
}

func (req *RequestParams) SetMethod(method string) *RequestParams {
	req.method = strings.ToUpper(method)
	return req
}

func (req *RequestParams) Get() *RequestParams {
	req.method = http.MethodGet
	return req
}

func (req *RequestParams) Post() *RequestParams {
	req.method = http.MethodPost
	return req
}

func (req *RequestParams) Put() *RequestParams {
	req.method = http.MethodPut
	return req
}

func (req *RequestParams) Delete() *RequestParams {
	req.method = http.MethodDelete
	return req
}

func (req *RequestParams) SetContentType(contentType ContentType) *RequestParams {
	req.ContentType = contentType
	return req
}

func (req *RequestParams) AddHeader(k, v string) *RequestParams {
	req.header.Set(k, v)
	return req
}

func (req *RequestParams) SetHeader(header http.Header) *RequestParams {
	req.header = header
	return req
}

func (req *RequestParams) CachedHeader(key string) *RequestParams {
	req.cachedHeaderKey = key
	return req
}

func (req *RequestParams) SetLogger(logger LogCallback) *RequestParams {
	req.logger = logger
	return req
}

func (req *RequestParams) DisableLog() *RequestParams {
	req.logger = nil
	return req
}

func (req *RequestParams) GenLog() *RequestParams {
	if req.logger == nil {
		req.logger = DefaultLogger
	}
	req.genlog = true
	return req
}

func (req *RequestParams) SetResponseHandler(handler func([]byte) ([]byte, error)) *RequestParams {
	req.ResponseHandler = handler
	return req
}

func (req *RequestParams) SetTimeout(timeout time.Duration) *RequestParams {
	req.timeout = timeout
	return req
}

func (req *RequestParams) SetClient(client *http.Client) *RequestParams {
	req.client = client
	return req
}

func (req *RequestParams) SetRetryTimes(retryTimes int) *RequestParams {
	req.retryTimes = retryTimes
	return req
}

func (req *RequestParams) SetRetryHandle(handle func(*RequestParams)) *RequestParams {
	req.retryHandle = handle
	return req
}

type ResponseBodyCheck interface {
	CheckError() error
}

type ResponseBody struct {
	Status  int         `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func CommonResponse(response interface{}) ResponseBodyCheck {
	return &ResponseBody{Data: response}
}

func (res *ResponseBody) CheckError() error {
	if res.Status != 0 {
		return errors.New(res.Message)
	}
	return nil
}

type RawResponse = []byte

// Do create a HTTP request
func (req *RequestParams) Do(param, response interface{}) error {
	method := req.method
	url := req.url
	if req.timeout != 0 {
		defer setTimeout(req.client, timeout)
		setTimeout(req.client, req.timeout)
	}
	var body io.Reader
	var reqBody, respBody *Body
	var statusCode int
	var err error
	reqTime := time.Now()
	// 日志记录
	defer func(now time.Time) {
		if req.genlog || (req.logger != nil && genlog) {
			req.logger(url, method, req.AuthUser, reqBody, respBody, statusCode, time.Now().Sub(now), err)
		}
	}(reqTime)

	if method == http.MethodGet {
		if param != nil {
			switch paramt := param.(type) {
			case string:
				url += "?" + paramt
			case []byte:
				url += "?" + stringsi.ToString(paramt)
			default:
				params := getParam(param)
				url += "?" + params
			}
		}
	} else {
		reqBody = &Body{}
		if param != nil {
			switch paramt := param.(type) {
			case string:
				body = strings.NewReader(paramt)
				reqBody.Data = stringsi.ToBytes(paramt)
			case []byte:
				body = bytes.NewReader(paramt)
				reqBody.Data = paramt
			case io.Reader:
				var reqBytes []byte
				reqBytes, err = io.ReadAll(paramt)
				body = bytes.NewReader(reqBytes)
				reqBody.Data = reqBytes
			default:
				if req.ContentType == ContentTypeJson {
					var reqBytes []byte
					reqBytes, err = json.Marshal(param)
					if err != nil {
						return err
					}
					body = bytes.NewReader(reqBytes)
					reqBody.Data = reqBytes
					reqBody.ContentType = ContentTypeJson
				} else {
					params := getParam(param)
					reqBody.Data = stringsi.ToBytes(params)
					body = strings.NewReader(params)
				}
			}
		}
	}
	var request *http.Request
	request, err = http.NewRequestWithContext(req.ctx, method, url, body)
	if err != nil {
		return err
	}

	// 缓存header
	if req.cachedHeaderKey != "" {
		if header, ok := headerMap.Load(req.cachedHeaderKey); ok {
			request.Header = header.(http.Header)
		}
	} else {
		request.Header = req.header
	}

	if req.AuthUser != "" && req.AuthPass != "" {
		request.SetBasicAuth(req.AuthUser, req.AuthPass)
	}
	if req.ContentType == ContentTypeJson {
		request.Header.Set(httpi.HeaderContentType, httpi.ContentJSONHeaderValue)
	} else if req.ContentType == ContentTypeFormData {
		request.Header.Set(httpi.HeaderContentType, httpi.ContentFormMultipartHeaderValue)
	} else {
		request.Header.Set(httpi.HeaderContentType, httpi.ContentFormHeaderValue)
	}
	var resp *http.Response
	resp, err = req.client.Do(request)
	if err != nil {
		if req.retryTimes == 0 {
			return err
		}
		for i := 0; i < req.retryTimes; i++ {
			if req.retryHandle != nil {
				req.retryHandle(req)
			}
			reqTime = time.Now()
			resp, err = req.client.Do(request)
			if err == nil {
				break
			} else {
				if req.genlog || (req.logger != nil && genlog) {
					req.logger(url, method, req.AuthUser, reqBody, respBody, statusCode, time.Now().Sub(reqTime), errors.New(err.Error()+";will retry"))
				}
			}
		}
	}
	respBody = &Body{}
	var respBytes []byte
	respBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = resp.Body.Close()
	if req.ResponseHandler != nil {
		respBytes, err = req.ResponseHandler(respBytes)
		if err != nil {
			return err
		}
	}
	respBody.Data = respBytes
	statusCode = resp.StatusCode
	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		respBody.ContentType = ContentTypeForm
		err = errors.New("status:" + resp.Status + "" + stringsi.ToString(respBytes))
		return err
	}
	if len(respBytes) > 0 && response != nil {
		if raw, ok := response.(*RawResponse); ok {
			*raw = respBytes
			if resp.Header.Get(httpi.HeaderContentType) == httpi.ContentFormHeaderValue {
				respBody.ContentType = ContentTypeForm
			} else {
				respBody.ContentType = ContentTypeJson
			}
			return nil
		}
		if resp.Header.Get(httpi.HeaderContentType) == httpi.ContentFormHeaderValue {
			// TODO
			respBody.ContentType = ContentTypeForm
		} else {
			// 默认json
			respBody.ContentType = ContentTypeJson
			err = json.Unmarshal(respBytes, response)
			if err != nil {
				return err
			}
		}

		if v, ok := response.(ResponseBodyCheck); ok {
			err = v.CheckError()
		}
	}

	return err
}

func (req *RequestParams) DoRaw(param, response interface{}) (RawResponse, error) {
	var raw RawResponse
	err := req.Do(param, &raw)
	if err != nil {
		return raw, err
	}
	return raw, json.Unmarshal(raw, response)
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

func Get(url string, response any) error {
	return NewGetRequest(url).Do(nil, response)
}

func (req *RequestParams) DoGet(url string, param, response interface{}) error {
	req.url = url
	req.method = http.MethodGet
	return req.Do(param, response)
}

func Post(url string) *RequestParams {
	return NewRequest(url, http.MethodPost)
}

func (req *RequestParams) DoPost(url string, param, response interface{}) error {
	req.url = url
	req.method = http.MethodPost
	return (req).Do(param, response)
}

func Put(url string) *RequestParams {
	return NewRequest(url, http.MethodPut)
}

func (req *RequestParams) DoPut(url string, param, response interface{}) error {
	req.url = url
	req.method = http.MethodPut
	return req.Do(param, response)
}

func Delete(url string) *RequestParams {
	return NewRequest(url, http.MethodDelete)
}

func (req *RequestParams) DoDelete(url string, param, response interface{}) error {
	req.url = url
	req.method = http.MethodDelete
	return req.Do(param, response)
}

func (req *RequestParams) CompleteDo(url, method string, param, response interface{}) error {
	req.url = url
	req.method = method
	return req.Do(param, response)
}

func (req *RequestParams) Download(url, path string) error {
	req.url = url
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := req.client.Do(request)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("错误码： %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

type ReplaceHttpRequest http.Request

func NewReplaceHttpRequest(r *http.Request) *ReplaceHttpRequest {
	return (*ReplaceHttpRequest)(r)
}

func (r *ReplaceHttpRequest) SetURL(url string) *ReplaceHttpRequest {
	u, err := urlpkg.Parse(url)
	if err != nil {
		log.Error(err)
	}
	u.Host = removeEmptyPort(u.Host)
	r.URL = u
	r.Host = u.Host
	return r
}

// Given a string of the form "host", "host:port", or "[ipv6::address]:port",
// return true if the string includes a port.
func hasPort(s string) bool { return strings.LastIndex(s, ":") > strings.LastIndex(s, "]") }

// removeEmptyPort strips the empty port in ":port" to ""
// as mandated by RFC 3986 Section 6.2.3.
func removeEmptyPort(host string) string {
	if hasPort(host) {
		return strings.TrimSuffix(host, ":")
	}
	return host
}

func (r *ReplaceHttpRequest) SetMethod(method string) *ReplaceHttpRequest {
	r.Method = strings.ToUpper(method)
	return r
}

func (r *ReplaceHttpRequest) SetBody(body io.ReadCloser) *ReplaceHttpRequest {
	r.Body = body
	return r
}

func (r *ReplaceHttpRequest) SetContext(ctx context.Context) *ReplaceHttpRequest {
	stdr := (*http.Request)(r).WithContext(ctx)
	return (*ReplaceHttpRequest)(stdr)
}

func (r *ReplaceHttpRequest) StdHttpRequest() *http.Request {
	return (*http.Request)(r)
}
