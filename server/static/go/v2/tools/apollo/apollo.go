package apollo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	url2 "net/url"

	"github.com/liov/hoper/go/v2/utils/log"
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
	url := fmt.Sprintf("http://%s/configfiles/json/%s/%s/%s", c.Addr, c.AppId, c.Cluster, namespace)
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

type Server struct {
	Addr       string
	AppId      string `json:"appId"`
	Cluster    string `json:"cluster"`
	IP         string `json:"ip"`
	InitConfig SpecialConfig
	//一个应用的appId和集群应该是启动时确定的，但却可以有多个namespace
	Configurations map[string]*Apollo
	Notifications  NotificationMap
	close          chan struct{}
}

type SpecialConfig struct {
	NameSpace string
	Callback  func(map[string]string)
}

func New(addr, appId, cluster, ip string, namespaces []string, initConfig SpecialConfig) *Server {
	s := &Server{
		Addr:           addr,
		AppId:          appId,
		Cluster:        cluster,
		IP:             ip,
		InitConfig:     initConfig,
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

func (s *Server) GetCacheConfig(namespace string) error {
	url := fmt.Sprintf("http://%s/configfiles/json/%s/%s/%s", s.Addr, s.AppId, s.Cluster, namespace)
	if s.IP != "" {
		url += "?ip=" + s.IP
	}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	_, ok := s.Configurations[namespace]
	if !ok {
		s.Configurations[namespace] = &Apollo{}
	}
	err = json.Unmarshal(body, &(s.Configurations[namespace].Configurations))
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) GetNoCacheConfig(namespace string) error {
	url := fmt.Sprintf("http://%s/configs/%s/%s/%s", s.Addr, s.AppId, s.Cluster, namespace)
	_, ok := s.Configurations[namespace]
	if !ok {
		s.Configurations[namespace] = &Apollo{}
	}
	releaseKey := s.Configurations[namespace].ReleaseKey
	if releaseKey != "" && s.IP != "" {
		url += "?releaseKey=" + releaseKey + "&ip=" + s.IP
	} else {
		if s.IP != "" {
			url += "?ip=" + s.IP
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
	err = json.Unmarshal(body, s.Configurations[namespace])
	if namespace == s.InitConfig.NameSpace {
		s.InitConfig.Callback(s.Configurations[namespace].Configurations)
		s.Configurations[namespace] = nil
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

func (s *Server) UpdateConfig(appId, clusterName string) error {
	urlStr := fmt.Sprintf("http://%s/notifications/v2?appId=%s&cluster=%s",
		s.Addr, appId, clusterName)
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
			notificationSlice := s.Notifications.ToSlice()
			notifications, err := json.Marshal(&notificationSlice)
			newQuery := url2.Values{
				"notifications": []string{string(notifications)},
			}.Encode()
			req.URL.RawQuery = query + newQuery
			log.Debug("发送请求", req.URL)
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
				s.Notifications[changeNamespace[i].NamespaceName] = changeNamespace[i].NotificationId
				err = s.GetNoCacheConfig(changeNamespace[i].NamespaceName)
				log.Info(s.Configurations[changeNamespace[i].NamespaceName])
				if err != nil {
					return err
				}
			}
			ch <- struct{}{}
		case <-s.close:
			req.Header["Connection"] = []string{"Closed"}
			http.DefaultClient.Do(req)
			break Loop
		}
	}
}

func (s *Server) Get(namespace, key string) string {
	if s.Configurations == nil {
		return ""
	}
	ap, ok := s.Configurations[namespace]
	if !ok {
		return ""
	}
	if ap.Configurations == nil {
		return ""
	}
	return ap.Configurations[key]
}

func (s *Server) GetDefault(key string) string {
	return s.Get("default", key)
}

func (s *Server) Close() error {
	s.close <- struct{}{}
}
