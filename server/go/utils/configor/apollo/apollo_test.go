package apollo

import (
	"testing"

	"github.com/liov/hoper/v2/utils/log"
)

func TestApollo(t *testing.T) {
	s := Client{
		Addr:           "192.168.1.212:8080",
		AppId:          "hoper",
		Cluster:        "default",
		IP:             "",
		Configurations: map[string]*Apollo{},
		Notifications:  map[string]int{"application": -1},
	}
	err := s.GetCacheConfig("application")
	if err != nil {
		log.Error(err)
	}
	log.Debug(s.Configurations["application"])
	err = s.GetNoCacheConfig("application")
	if err != nil {
		log.Error(err)
	}
	log.Info(s.Configurations["application"])
}

func TestApolloPoll(t *testing.T) {
	s := Client{
		Addr:           "192.168.1.212:8080",
		AppId:          "hoper",
		Cluster:        "default",
		IP:             "",
		Configurations: map[string]*Apollo{},
		Notifications:  map[string]int{"application": -1},
	}
	err := s.UpdateConfig("hoper", "default")
	if err != nil {
		log.Error(err)
	}
}
