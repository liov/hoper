package apollo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	url2 "net/url"

	"github.com/liov/hoper/v2/utils/log"
)

const (
	GetConfigUrl     = "http://%s/configfiles/json/%s/%s/%s"
	NoCacheUrl       = "http://%s/configs/%s/%s/%s"
	NotificationsUrl = "http://%s/notifications/v2?appId=%s&cluster=%s"
)

type Config struct {
	Addr    string
	AppId   string `json:"appId"`
	Cluster string `json:"cluster"`
	IP      string `json:"ip"`
}

func NewConfig(addr, appId, cluster, ip string) *Config {
	return &Config{
		Addr:    addr,
		AppId:   appId,
		Cluster: cluster,
		IP:      ip,
	}
}

func (c *Config) GetInitConfig(namespace string) (map[string]string, error) {
	url := fmt.Sprintf(GetConfigUrl, c.Addr, c.AppId, c.Cluster, namespace)
	if c.IP != "" {
		url += "?ip=" + c.IP
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var config map[string]string
	err = json.Unmarshal(body, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

type Client struct {
	Addr          string
	AppId         string `json:"appId"`
	Cluster       string `json:"cluster"`
	IP            string `json:"ip"`
	SpecialConfig SpecialConfig
	//一个应用的appId和集群应该是启动时确定的，但却可以有多个namespace
	Configurations map[string]*Apollo
	Notifications  NotificationMap
	close          chan struct{}
}

type SpecialConfig struct {
	NameSpace string
	Callback  func(map[string]string)
}

func New(addr, appId, cluster, ip string, namespaces []string, initConfig SpecialConfig) *Client {
	s := &Client{
		Addr:           addr,
		AppId:          appId,
		Cluster:        cluster,
		IP:             ip,
		SpecialConfig:  initConfig,
		Configurations: map[string]*Apollo{},
		Notifications:  map[string]int{},
		close:          make(chan struct{}, 1),
	}

	for i := range namespaces {
		s.Notifications[namespaces[i]] = -1
		err := s.GetCacheConfig(namespaces[i])
		if err != nil {
			log.Error(err)
		}
	}

	go func() {
		s.UpdateConfig(s.AppId, s.Cluster)
	}()

	return s
}

type Apollo struct {
	AppId          string            `json:"appId"`
	Cluster        string            `json:"cluster"`
	NamespaceName  string            `json:"namespaceName"`
	Configurations map[string]string `json:"configurations"`
	ReleaseKey     string            `json:"releaseKey"`
}

type NoChangeError struct{}

func (e *NoChangeError) Error() string {
	return "配置无变化"
}

func (c *Client) GetCacheConfig(namespace string) error {
	url := fmt.Sprintf(GetConfigUrl, c.Addr, c.AppId, c.Cluster, namespace)
	if c.IP != "" {
		url += "?ip=" + c.IP
	}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	_, ok := c.Configurations[namespace]

	if !ok {
		c.Configurations[namespace] = &Apollo{}
	}
	err = json.Unmarshal(body, &(c.Configurations[namespace].Configurations))
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetNoCacheConfig(namespace string) error {
	url := fmt.Sprintf(NoCacheUrl, c.Addr, c.AppId, c.Cluster, namespace)
	_, ok := c.Configurations[namespace]
	if !ok {
		c.Configurations[namespace] = &Apollo{}
	}
	releaseKey := c.Configurations[namespace].ReleaseKey
	if releaseKey != "" && c.IP != "" {
		url += "?releaseKey=" + releaseKey + "&ip=" + c.IP
	} else {
		if c.IP != "" {
			url += "?ip=" + c.IP
		}
		if releaseKey != "" {
			url += "?releaseKey=" + releaseKey
		}
	}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, c.Configurations[namespace])
	if namespace == c.SpecialConfig.NameSpace {
		c.SpecialConfig.Callback(c.Configurations[namespace].Configurations)
		c.Configurations[namespace].Configurations = nil
	}
	if err != nil {
		return err
	}
	return nil
}

type NotificationInfo struct {
	NamespaceName  string `json:"namespaceName"`
	NotificationId int    `json:"notificationId"`
}

type NotificationSlice []NotificationInfo

func (ns NotificationSlice) ToMap() NotificationMap {
	var nm = NotificationMap{}
	for i := range ns {
		nm[ns[i].NamespaceName] = ns[i].NotificationId
	}
	return nm
}

func (ns *NotificationSlice) Add(nameSpace string) {
	for i := range *ns {
		if (*ns)[i].NamespaceName == nameSpace {
			return
		}
	}
	*ns = append(*ns, NotificationInfo{nameSpace, -1})
}

type NotificationMap map[string]int

func (nm NotificationMap) ToSlice() NotificationSlice {
	var ns NotificationSlice
	for k, v := range nm {
		ns = append(ns, NotificationInfo{k, v})
	}
	return ns
}

func (nm NotificationMap) Update(ns []NotificationInfo) {
	for i := range ns {
		nm[ns[i].NamespaceName] = ns[i].NotificationId
	}
}

func (c *Client) UpdateConfig(appId, clusterName string) error {
	urlStr := fmt.Sprintf(NotificationsUrl,
		c.Addr, appId, clusterName)
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return err
	}
	req.Header["Connection"] = []string{"keep-alive"}
	query := req.URL.RawQuery + "&"
	var ch = make(chan struct{}, 1)
	ch <- struct{}{}
Loop:
	for {
		select {
		case <-ch:
			notificationSlice := c.Notifications.ToSlice()
			notifications, err := json.Marshal(&notificationSlice)
			newQuery := url2.Values{
				"notifications": []string{string(notifications)},
			}.Encode()
			req.URL.RawQuery = query + newQuery
			log.Debug("发送请求:", req.URL)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return err
			}
			if resp.StatusCode == 304 {
				ch <- struct{}{}
				continue
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			log.Debug("收到返回", string(body))
			var changeNamespace []NotificationInfo
			err = json.Unmarshal(body, &changeNamespace)
			if err != nil {
				return err
			}
			for i := range changeNamespace {
				c.Notifications[changeNamespace[i].NamespaceName] = changeNamespace[i].NotificationId
				err = c.GetNoCacheConfig(changeNamespace[i].NamespaceName)
				log.Info(c.Configurations[changeNamespace[i].NamespaceName])
				if err != nil {
					return err
				}
			}
			ch <- struct{}{}
		case <-c.close:
			req.Header["Connection"] = []string{"Closed"}
			http.DefaultClient.Do(req)
			break Loop
		}
	}
	return nil
}

func (c *Client) Get(namespace, key string) string {
	if c.Configurations == nil {
		return ""
	}
	ap, ok := c.Configurations[namespace]
	if !ok {
		return ""
	}
	if ap.Configurations == nil {
		return ""
	}
	return ap.Configurations[key]
}

func (c *Client) GetDefault(key string) string {
	return c.Get("default", key)
}

func (c *Client) Close() error {
	c.close <- struct{}{}
	return nil
}
