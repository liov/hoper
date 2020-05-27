package cache

import (
	"net/http"

	"github.com/golang/groupcache"
	pb "github.com/golang/groupcache/groupcachepb"
	"github.com/liov/hoper/go/v2/utils/log"
)

type Server struct {
	Port       string
	BasePath   string
	GroupName  string
	CacheBytes int64
}

func (s *Server) ServerHttp() {
	me := "http://localhost" + s.Port
	opts := groupcache.HTTPPoolOptions{BasePath: s.BasePath}
	peers := groupcache.NewHTTPPoolOpts(me, &opts)

	peers.Set("http://localhost:8333", "http://localhost:8222")
	groupcache.NewGroup(s.GroupName, s.CacheBytes, groupcache.GetterFunc(
		func(ctx groupcache.Context, key string, dest groupcache.Sink) error {
			dest.SetString(key + me)
			return nil
		}))

	http.HandleFunc(s.BasePath, peers.ServeHTTP)
	http.ListenAndServe(s.Port, nil)
}

func GetFromPeer(groupName, key string, peers *groupcache.HTTPPool) (value []byte, err error) {
	req := &pb.GetRequest{Group: &groupName, Key: &key}
	res := &pb.GetResponse{}

	peer, ok := peers.PickPeer(key)
	if ok == false {
		log.Info("peers PickPeer failed: ", key)
		return
	}

	err = peer.Get(nil, req, res)
	if err != nil {
		return nil, err
	}
	return res.Value, nil
}
