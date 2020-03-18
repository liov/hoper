package nacos

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/liov/hoper/go/v2/utils/log"
)

const (
	GetConfigUrl        = "http://%s/nacos/v1/cs/configs?tenant=%s&group=%s&dataId=%s"
	GetConfigAllInfoUrl = "http://%s/nacos/v1/cs/configs?show=all&tenant=%s&group=%s&dataId=%s"
	ListenerUrl         = "http://%s/nacos/v1/cs/configs/listener"
	InitParam           = "Listening-Configs=%s" + string(rune(2)) + "%s" + string(rune(2)) + "%s" + string(rune(2)) + "%s" + string(rune(1))
)

type Config struct {
	Addr   string `json:"addr"`
	Tenant string `json:"tenant"`
	Group  string `json:"group"`
	DataId string `json:"dataId"`
}

type Client struct {
	*Config
	MD5   string `json:"md5"`
	close chan struct{}
}

func (c *Config) NewClient() *Client {
	return &Client{
		Config: c,
		close:  make(chan struct{}),
	}
}

type ConfigInfo struct {
	ID         int64  `json:"id"`
	AppName    string `json:"appName"`
	Content    string `json:"content"`
	CreateIp   string `json:"createIp"`
	CreateTime int64  `json:"createTime"`
	DataId     string `json:"dataId"`
	Desc       string `json:"desc"`
	Effect     string `json:"effect"`
	Group      string `json:"group"`
	MD5        string `json:"md5"`
	ModifyTime int64  `json:"modifyTime"`
	Type       string `json:"type"`
}

func (c *Client) GetConfig() ([]byte, error) {
	urlStr := fmt.Sprintf(GetConfigUrl,
		c.Addr, c.Tenant, c.Group, c.DataId)
	resp, err := http.Get(urlStr)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}

func (c *Client) GetConfigAllInfo() ([]byte, error) {
	urlStr := fmt.Sprintf(GetConfigAllInfoUrl,
		c.Addr, c.Tenant, c.Group, c.DataId)
	resp, err := http.Get(urlStr)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var configInfo ConfigInfo
	json.Unmarshal(data, &configInfo)
	c.MD5 = configInfo.MD5
	return []byte(configInfo.Content), nil
}

func (c *Client) GetConfigAllInfoHandle(handle func([]byte)) error {
	config, err := c.GetConfigAllInfo()
	if err != nil {
		return err
	}
	handle(config)
	return nil
}

func (c *Client) Listener(handle func([]byte)) {
	urlStr := fmt.Sprintf(ListenerUrl, c.Addr)
	listeningConfigs := fmt.Sprintf(InitParam, c.DataId, c.Group, c.MD5, c.Tenant)
	req, err := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(listeningConfigs))
	if err != nil {
		log.Fatal(err)
	}
	req.Header["Long-Pulling-Timeout"] = []string{"30000"}
	req.Header["Content-Type"] = []string{"application/x-www-form-urlencoded"}

	var ch = make(chan struct{}, 1)
	ch <- struct{}{}
Loop:
	for {
		select {
		case <-ch:
			listeningConfigs = fmt.Sprintf(InitParam, c.DataId, c.Group, c.MD5, c.Tenant)
			req.Body = ioutil.NopCloser(strings.NewReader(listeningConfigs))
			log.Debug("发送请求", req.URL)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Error(err)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Error(err)
			}
			if resp.StatusCode != 200 {
				log.Error(string(body))
				ch <- struct{}{}
				continue
			}

			if len(body) != 0 {
				params := strings.Split(string(body), "%02")
				c.MD5 = params[2][:len(params[2])-4]
				c.GetConfigAllInfoHandle(handle)
			}
			ch <- struct{}{}
		case <-c.close:
			req.Header["Connection"] = []string{"Closed"}
			http.DefaultClient.Do(req)
			break Loop
		}
	}
}
