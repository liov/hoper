package apollo

import (
	"log"
	"testing"
)

func TestApollo(t *testing.T) {
	log.SetFlags(15)
	s := Server{
		Addr:           "192.168.1.212:8080",
		AppId:          "hoper",
		Cluster:        "default",
		IP:             "",
		Configurations: map[string]*Apollo{},
		Notifications:  map[string]int{"application": -1},
	}
	err := s.GetCacheConfig("application")
	if err != nil {
		log.Println(err)
	}
	log.Println(s.Configurations["application"])
	err = s.GetNoCacheConfig("application")
	if err != nil {
		log.Println(err)
	}
	log.Println(s.Configurations["application"])
}
