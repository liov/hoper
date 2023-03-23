package client

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	httpi "github.com/liov/hoper/server/go/lib/utils/net/http"
	"github.com/liov/hoper/server/go/lib/utils/strings"
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
	ContentTypeImage
)

// Request ...
type Request struct {
	ctx                context.Context
	client             *http.Client
	url, method        string
	contentType        ContentType
	timeout            time.Duration
	AuthUser, AuthPass string
	header             Header
	cachedHeaderKey    string
	logger             LogCallback
	disableLog         bool
	responseHandler    func(response []byte) ([]byte, error)
	retryTimes         int
	retryHandle        func(*Request)
}

func New(url string) *Request {
	return newRequest(url, "")
}

func NewRequest(url, method string) *Request {
	return newRequest(url, strings.ToUpper(method))
}

func newRequest(url, method string) *Request {
	return &Request{ctx: context.Background(), client: defaultClient, url: url, method: method, header: make([]string, 0, 2), logger: defaultLog}
}

func NewGetRequest(url string) *Request {
	return newRequest(url, http.MethodGet)
}

func NewPostRequest(url string) *Request {
	return newRequest(url, http.MethodPost)
}

func NewPutRequest(url string) *Request {
	return newRequest(url, http.MethodPut)
}

func NewDeleteRequest(url string) *Request {
	return newRequest(url, http.MethodDelete)
}

func (req *Request) Url(url string) *Request {
	req.url = url
	return req
}

func (req *Request) Method(method string) *Request {
	req.method = strings.ToUpper(method)
	return req
}

func (req *Request) ContentType(contentType ContentType) *Request {
	req.contentType = contentType
	return req
}

func (req *Request) SetHeader(header Header) *Request {
	req.header = header
	return req
}

func (req *Request) AddHeader(k, v string) *Request {
	req.header = append(req.header, k, v)
	return req
}

func (req *Request) CachedHeader(key string) *Request {
	req.cachedHeaderKey = key
	return req
}

func (req *Request) WithLogger(logger LogCallback) *Request {
	req.logger = logger
	return req
}

func (req *Request) DisableLog() *Request {
	req.disableLog = true
	return req
}

func (req *Request) ResponseHandler(handler func([]byte) ([]byte, error)) *Request {
	req.responseHandler = handler
	return req
}

func (req *Request) Timeout(timeout time.Duration) *Request {
	req.timeout = timeout
	return req
}

func (req *Request) WithClient(client *http.Client) *Request {
	req.client = client
	return req
}

func (req *Request) RetryTimes(retryTimes int) *Request {
	req.retryTimes = retryTimes
	return req
}

func (req *Request) RetryHandle(handle func(*Request)) *Request {
	req.retryHandle = handle
	return req
}

func (req *Request) UrlParam(param interface{}) *Request {
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

func (req *Request) DoWithNoParam(response interface{}) error {
	return req.Do(nil, response)
}

func (req *Request) DoWithNoResponse(param interface{}) error {
	return req.Do(param, nil)
}

func (req *Request) DoEmpty() error {
	return req.Do(nil, nil)
}

func (req *Request) setHeader(request *http.Request) {
	for i := 0; i+1 < len(req.header); i += 2 {
		request.Header.Set(req.header[i], req.header[i+1])
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
}

// Do create a HTTP request
func (req *Request) Do(param, response interface{}) error {
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
		if !req.disableLog {
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
		} else {
			req.setHeader(request)
			headerMap.Store(req.cachedHeaderKey, request.Header)
		}
	} else {
		req.setHeader(request)
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
			time.Sleep(time.Millisecond * 200)
			reqTime = time.Now()
			request.Body = io.NopCloser(bytes.NewReader(reqBody.Data))
			resp, err = req.client.Do(request)
			if err == nil {
				break
			} else if !req.disableLog {
				req.logger(url, method, req.AuthUser, reqBody, respBody, statusCode, time.Since(reqTime), errors.New(err.Error()+";will retry"))
			}
		}
		if err != nil {
			return err
		}
	}
	if httpresp, ok := response.(**http.Response); ok {
		*httpresp = resp
		return err
	}

	var reader io.Reader
	if resp.Header.Get(httpi.HeaderContentEncoding) == "gzip" {
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			if resp != nil {
				resp.Body.Close()
			}
			return err
		}
	} else {
		reader = resp.Body
	}
	if httpresp, ok := response.(*io.Reader); ok {
		*httpresp = reader
		return err
	}
	respBody = &Body{}
	var respBytes []byte
	respBytes, err = io.ReadAll(reader)
	if err != nil {
		return err
	}
	resp.Body.Close()
	statusCode = resp.StatusCode
	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		respBody.ContentType = ContentTypeText
		if resp.StatusCode == http.StatusNotFound {
			err = errors.New("not found")
		} else {
			err = errors.New("status:" + resp.Status + " " + stringsi.ConvertUnicode(respBytes))
		}
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
		contentType := resp.Header.Get(httpi.HeaderContentType)
		if contentType == httpi.ContentJSONHeaderValue {
			respBody.ContentType = ContentTypeJson
		} else if contentType == httpi.ContentFormHeaderValue {
			respBody.ContentType = ContentTypeForm
		} else if strings.HasPrefix(contentType, "text") {
			respBody.ContentType = ContentTypeText
		} else if strings.HasPrefix(contentType, "image") {
			respBody.ContentType = ContentTypeImage
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
				return errors.New("json.Unmarshal error:" + err.Error() + " status:" + resp.Status + " " + stringsi.ConvertUnicode(respBytes))
			}
		}

		if v, ok := response.(ResponseBodyCheck); ok {
			err = v.CheckError()
		}
	}

	return err
}

func (req *Request) DoRaw(param interface{}) (RawResponse, error) {
	var raw RawResponse
	err := req.Do(param, &raw)
	if err != nil {
		return raw, err
	}
	return raw, nil
}

func (req *Request) DoStream(param interface{}) (io.ReadCloser, error) {
	var resp *http.Response
	err := req.Do(param, &resp)
	if err != nil {
		return resp.Body, err
	}
	return resp.Body, nil
}

func (req *Request) Get(url string, response interface{}) error {
	req.url = url
	req.method = http.MethodGet
	return req.Do(nil, response)
}

func (req *Request) Post(url string, param, response interface{}) error {
	req.url = url
	req.method = http.MethodPost
	return (req).Do(param, response)
}

func (req *Request) Put(url string, param, response interface{}) error {
	req.url = url
	req.method = http.MethodPut
	return req.Do(param, response)
}

func (req *Request) Delete(url string, param, response interface{}) error {
	req.url = url
	req.method = http.MethodDelete
	return req.Do(param, response)
}

func (req *Request) CompleteDo(url, method string, param, response interface{}) error {
	req.url = url
	req.method = method
	return req.Do(param, response)
}

func (req *Request) Download(url, path string) error {
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

func (req *Request) CacheDo(url, method string, param, response interface{}) error {
	req.url = url
	req.method = method
	return req.Do(param, response)
}

func (req *Request) CacheGet(url string, response interface{}) error {
	req.url = url
	req.method = http.MethodGet
	return req.Do(nil, response)
}

func (req *Request) CachePost(url string, param, response interface{}) error {
	req.url = url
	req.method = http.MethodPost
	return req.Do(param, response)
}

func (req *Request) CachePut(url string, param, response interface{}) error {
	req.url = url
	req.method = http.MethodPut
	return req.Do(param, response)
}

func (req *Request) CacheDelete(url string, param, response interface{}) error {
	req.url = url
	req.method = http.MethodDelete
	return req.Do(param, response)
}

type SetParams func(req *Request) *Request
