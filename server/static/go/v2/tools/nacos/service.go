package nacos

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kataras/iris/v12/context"
)

const (
	CreateServiceUrl  = "http://%s/nacos/v1/ns/service?serviceName=%s&groupName=%s&namespaceId=%s"
	GetServiceListUrl = "http://%s/nacos/v1/ns/service/list?pageNo=%s&pageSize=%s&serviceName=%s&groupName=%s&namespaceId=%s"
)

func (c *Client) CreateService(svcName string) error {
	urlStr := fmt.Sprintf(CreateServiceUrl,
		c.Addr, svcName, c.Group, c.Tenant)
	resp, err := http.Post(urlStr, context.ContentFormHeaderValue, nil)
	if err != nil {
		return err
	}
	if res, _ := ioutil.ReadAll(resp.Body); string(res) != "ok" {
		return errors.New(fmt.Sprintf("服务:%s注册失败", svcName))
	}
	return nil
}

func (c *Client) GetService(svcName string) (string, error) {
	urlStr := fmt.Sprintf(CreateServiceUrl,
		c.Addr, svcName, c.Group, c.Tenant)
	resp, err := http.Get(urlStr)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("服务:%s获取失败", svcName))
	}
	res, _ := ioutil.ReadAll(resp.Body)
	if string(res) != "ok" {
		return "", errors.New(fmt.Sprintf("服务:%s获取失败", svcName))
	}
	return string(res), nil
}
