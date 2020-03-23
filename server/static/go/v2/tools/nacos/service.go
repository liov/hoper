package nacos

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/kataras/iris/v12/context"
	"github.com/liov/hoper/go/v2/utils/net"
)

const (
	CreateServiceUrl      = "http://%s/nacos/v1/ns/service"
	CreateServiceParam    = "serviceName=%s&groupName=%s&namespaceId=%s&protectThreshold=0&metadata=%s"
	GetServiceUrl         = "http://%s/nacos/v1/ns/service?serviceName=%s&groupName=%s&namespaceId=%s&protectThreshold=0"
	GetServiceListUrl     = "http://%s/nacos/v1/ns/service/list?pageNo=%s&pageSize=%s&serviceName=%s&groupName=%s&namespaceId=%s"
	RegisterInstanceUrl   = "http://%s/nacos/v1/ns/instance"
	RegisterInstanceParam = "port=%s&ip=%s&serviceName=%s&groupName=%s&namespaceId=%s"
	BeatParam             = "serviceName=%s&groupName=%s&beat=%s"
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

func (c *Config) CreateService(svcName string, metadata *Metadata) error {
	urlStr := fmt.Sprintf(CreateServiceUrl, c.Addr)
	var data []byte
	if metadata != nil {
		data, _ = json.Marshal(metadata)
	}
	param := fmt.Sprintf(CreateServiceParam, svcName, c.Group, c.Tenant, string(data))
	resp, err := http.Post(urlStr, context.ContentFormHeaderValue, strings.NewReader(param))
	if err != nil {
		return err
	}
	if res, _ := ioutil.ReadAll(resp.Body); resp.StatusCode != 200 {
		return errors.New(string(res))
	}
	return nil
}

func (c *Config) GetService(svcName string) (*Service, error) {
	urlStr := fmt.Sprintf(GetServiceUrl,
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

func (c *Config) RegisterInstance(port, svcName string) error {
	urlStr := fmt.Sprintf(RegisterInstanceUrl, c.Addr)
	param := fmt.Sprintf(RegisterInstanceParam,
		port[1:], net.GetIP(), svcName, c.Group, c.Tenant)
	resp, err := http.Post(urlStr, context.ContentFormHeaderValue, strings.NewReader(param))
	if err != nil {
		return err
	}
	if res, _ := ioutil.ReadAll(resp.Body); resp.StatusCode != 200 {
		return errors.New(string(res))
	}
	return nil
}

func (c *Config) DeleteInstance(port, svcName string) error {
	urlStr := fmt.Sprintf(RegisterInstanceUrl, c.Addr)
	param := fmt.Sprintf(RegisterInstanceParam,
		port[1:], net.GetIP(), svcName, c.Group, c.Tenant)
	req, err := http.NewRequest(http.MethodDelete, urlStr, strings.NewReader(param))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", context.ContentFormHeaderValue)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res, _ := ioutil.ReadAll(resp.Body); resp.StatusCode != 200 {
		return errors.New(string(res))
	}
	return nil
}

func (c *Config) InstanceBeat(svcName string) error {
	urlStr := fmt.Sprintf(RegisterInstanceUrl, c.Addr) + "/beat"
	param := fmt.Sprintf(BeatParam, svcName, c.Group, `{"msg":"实例正常"}`)
	req, err := http.NewRequest(http.MethodPut, urlStr, strings.NewReader(param))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", context.ContentFormHeaderValue)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res, _ := ioutil.ReadAll(resp.Body); resp.StatusCode != 200 {
		return errors.New(string(res))
	}
	return nil
}
