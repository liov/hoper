package apollo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	url2 "net/url"
)

type Server struct {
	Addr    string
	AppId   string `json:"appId"`
	Cluster string `json:"cluster"`
	IP      string `json:"ip"`
	//一个应用的appId和集群应该是启动时确定的，但却可以有多个namespace
	Configurations map[string]*Apollo
	Notifications  NotificationMap
}

func New(addr, appId, cluster, ip string) *Server {
	return &Server{
		Addr:           addr,
		AppId:          appId,
		Cluster:        cluster,
		IP:             ip,
		Configurations: map[string]*Apollo{},
		Notifications:  map[string]int{},
	}
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

	_, ok := s.Configurations[namespace]
	if !ok {
		s.Configurations[namespace] = &Apollo{}
	}
	err = json.Unmarshal(body, s.Configurations[namespace])
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
	notificationSlice := s.Notifications.ToSlice()
	notifications, err := json.Marshal(&notificationSlice)
	url := fmt.Sprintf("http://%s/notifications/v2?", s.Addr) +
		url2.Values{
			"appId":         []string{appId},
			"cluster":       []string{clusterName},
			"notifications": []string{string(notifications)},
		}.Encode()
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if resp.Status == "304" {
		return &NoChangeError{}
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var changeNamespace []NotificationInfo
	err = json.Unmarshal(body, &changeNamespace)
	if err != nil {
		return err
	}
	for i := range changeNamespace {
		s.Notifications[changeNamespace[i].NamespaceName] = changeNamespace[i].NotificationId
		err = s.GetNoCacheConfig(changeNamespace[i].NamespaceName)
		if err != nil {
			return err
		}
	}
	return nil
}
