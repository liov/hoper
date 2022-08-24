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

	httpi "github.com/actliboy/hoper/server/go/lib/utils/net/http"
	"github.com/actliboy/hoper/server/go/lib/utils/number"
	"github.com/actliboy/hoper/server/go/lib/utils/strings"
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
		Proxy:             http.ProxyFromEnvironment, // 代理使用
		ForceAttemptHTTP2: true,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(network, addr, timeout)
			if err != nil {
				return nil, err
			}
			err = c.SetDeadline(time.Now().Add(timeout))
			return c, err
		},
		//DisableKeepAlives: true,
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
	client.Transport.(*http.Transport).DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		c, err := net.DialTimeout(network, addr, timeout)
		if err != nil {
			return nil, err
		}
		err = c.SetDeadline(time.Now().Add(timeout))
		return c, err
	}

}

type Pair struct {
	K, V string
}

type ContentType uint8

const (
	ContentTypeJson ContentType = iota
	ContentTypeForm
	ContentTypeFormData
	ContentTypeProtobuf
	ContentTypeText
)

// RequestParams ...
type RequestParams struct {
	ctx                context.Context
	client             *http.Client
	url, method        string
	contentType        ContentType
	timeout            time.Duration
	AuthUser, AuthPass string
	header             http.Header
	cachedHeaderKey    string
	logger             LogCallback
	genlog             bool
	responseHandler    func(response []byte) ([]byte, error)
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

func (req *RequestParams) Method(method string) *RequestParams {
	req.method = strings.ToUpper(method)
	return req
}

func (req *RequestParams) ContentType(contentType ContentType) *RequestParams {
	req.contentType = contentType
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

func (req *RequestParams) WithLogger(logger LogCallback) *RequestParams {
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

func (req *RequestParams) ResponseHandler(handler func([]byte) ([]byte, error)) *RequestParams {
	req.responseHandler = handler
	return req
}

func (req *RequestParams) Timeout(timeout time.Duration) *RequestParams {
	req.timeout = timeout
	return req
}

func (req *RequestParams) WithClient(client *http.Client) *RequestParams {
	req.client = client
	return req
}

func (req *RequestParams) RetryTimes(retryTimes int) *RequestParams {
	req.retryTimes = retryTimes
	return req
}

func (req *RequestParams) RetryHandle(handle func(*RequestParams)) *RequestParams {
	req.retryHandle = handle
	return req
}

func (req *RequestParams) UrlParam(param interface{}) *RequestParams {
	if param == nil {
		return req
	}
	sep := "?"
	if strings.Contains(req.url, sep) {
		sep = "&"
	}
	switch paramt := param.(type) {
	case string:
		req.url += sep + paramt
	case []byte:
		req.url += sep + stringsi.ToString(paramt)
	default:
		params := UrlParam(param)
		req.url += sep + params
	}
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

func (req *RequestParams) DoWithNoParam(response interface{}) error {
	return req.Do(nil, response)
}

func (req *RequestParams) DoWithNoResponse(param interface{}) error {
	return req.Do(param, nil)
}

func (req *RequestParams) DoEmpty() error {
	return req.Do(nil, nil)
}

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
			req.logger(url, method, req.AuthUser, reqBody, respBody, statusCode, time.Since(now), err)
		}
	}(reqTime)

	if method == http.MethodGet {
		url = req.UrlParam(param).url
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
				if req.contentType == ContentTypeJson {
					var reqBytes []byte
					reqBytes, err = json.Marshal(param)
					if err != nil {
						return err
					}
					body = bytes.NewReader(reqBytes)
					reqBody.Data = reqBytes
					reqBody.ContentType = ContentTypeJson
				} else {
					params := UrlParam(param)
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
	if req.contentType == ContentTypeJson {
		request.Header.Set(httpi.HeaderContentType, httpi.ContentJSONHeaderValue)
	} else if req.contentType == ContentTypeFormData {
		request.Header.Set(httpi.HeaderContentType, httpi.ContentFormHeaderValue)
	} else {
		request.Header.Set(httpi.HeaderContentType, httpi.ContentFormMultipartHeaderValue)
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
					req.logger(url, method, req.AuthUser, reqBody, respBody, statusCode, time.Since(reqTime), errors.New(err.Error()+";will retry"))
				}
			}
		}
	}
	if httpresp, ok := response.(**http.Response); ok {
		*httpresp = resp
		return err
	}
	if httpresp, ok := response.(*io.ReadCloser); ok {
		*httpresp = resp.Body
		return err
	}
	respBody = &Body{}
	var respBytes []byte
	respBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()
	statusCode = resp.StatusCode
	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		respBody.ContentType = ContentTypeText
		err = errors.New("status:" + resp.Status + "" + stringsi.ToString(respBytes))
		return err
	}

	if req.responseHandler != nil {
		respBytes, err = req.responseHandler(respBytes)
		if err != nil {
			return err
		}
	}
	respBody.Data = respBytes
	if len(respBytes) > 0 && response != nil {
		if resp.Header.Get(httpi.HeaderContentType) == httpi.ContentFormHeaderValue {
			respBody.ContentType = ContentTypeForm
		} else {
			respBody.ContentType = ContentTypeJson
		}

		if raw, ok := response.(*RawResponse); ok {
			*raw = respBytes
			return nil
		}
		if respBody.ContentType == ContentTypeForm {
			// TODO
		} else {
			// 默认json
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

func (req *RequestParams) DoRaw(param interface{}) (RawResponse, error) {
	var raw RawResponse
	err := req.Do(param, &raw)
	if err != nil {
		return raw, err
	}
	return raw, nil
}

func (req *RequestParams) DoStream(param interface{}) (io.ReadCloser, error) {
	var resp *http.Response
	err := req.Do(param, &resp)
	if err != nil {
		return resp.Body, err
	}
	return resp.Body, nil
}

func GetStream(url string) (io.ReadCloser, error) {
	var resp *http.Response
	err := Get(url, &resp)
	if err != nil {
		return resp.Body, err
	}
	return resp.Body, nil
}

func UrlParam(param interface{}) string {
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
	return NewGetRequest(url).DoWithNoParam(response)
}

func (req *RequestParams) Get(url string, response interface{}) error {
	req.url = url
	req.method = http.MethodGet
	return req.Do(nil, response)
}

func Post(url string, param, response interface{}) error {
	return NewPostRequest(url).Do(param, response)
}

func (req *RequestParams) Post(url string, param, response interface{}) error {
	req.url = url
	req.method = http.MethodPost
	return (req).Do(param, response)
}

func Put(url string, param, response interface{}) error {
	return NewPutRequest(url).Do(param, response)
}

func (req *RequestParams) Put(url string, param, response interface{}) error {
	req.url = url
	req.method = http.MethodPut
	return req.Do(param, response)
}

func Delete(url string, param, response interface{}) error {
	return NewDeleteRequest(url).Do(param, response)
}

func (req *RequestParams) Delete(url string, param, response interface{}) error {
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
