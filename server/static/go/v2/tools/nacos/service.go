package nacos

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kataras/iris/v12/context"
	"github.com/liov/hoper/go/v2/utils/net"
)

const (
	CreateServiceUrl    = "http://%s/nacos/v1/ns/service?serviceName=%s&groupName=%s&namespaceId=%s"
	GetServiceListUrl   = "http://%s/nacos/v1/ns/service/list?pageNo=%s&pageSize=%s&serviceName=%s&groupName=%s&namespaceId=%s"
	RegisterInstanceUrl = "http://%s/nacos/v1/ns/instance?port=%s&ip=%s&serviceName=%s&groupName=%s&namespaceId=%s"
)

type Type struct {
	Type string `json:"type"`
}

type Metadata struct {
	Domain string `json:"domain"`
}

type Service struct {
	Metadata         Metadata `json:"metadata"`
	Name             string   `json:"name"`
	Selector         Type     `json:"selector"`
	ProtectThreshold float32  `json:"protectThreshold"`
	Clusters         []struct {
		HealthChecker Type   `json:"healthChecker"`
		Name          string `json:"name"`
	}
}

func (c *Client) CreateService(svcName string) error {
	urlStr := fmt.Sprintf(CreateServiceUrl,
		c.Addr, svcName, c.Group, c.Tenant)
	resp, err := http.Post(urlStr, context.ContentFormHeaderValue, nil)
	if err != nil {
		return err
	}
	if res, _ := ioutil.ReadAll(resp.Body); resp.StatusCode != 200 {
		return errors.New(string(res))
	}
	return nil
}

func (c *Client) GetService(svcName string) (*Service, error) {
	urlStr := fmt.Sprintf(CreateServiceUrl,
		c.Addr, svcName, c.Group, c.Tenant)
	resp, err := http.Get(urlStr)
	if err != nil {
		return nil, err
	}
	res, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, errors.New(string(res))
	}
	var service Service
	json.Unmarshal(res, &service)
	return &service, nil
}

func (c *Client) RegisterInstance(port, svcName string) error {
	urlStr := fmt.Sprintf(RegisterInstanceUrl,
		c.Addr, port, net.GetIP(), svcName, c.Group, c.Tenant)
	resp, err := http.Post(urlStr, context.ContentFormHeaderValue, nil)
	if err != nil {
		return err
	}
	if res, _ := ioutil.ReadAll(resp.Body); resp.StatusCode != 200 {
		return errors.New(string(res))
	}
	return nil
}
