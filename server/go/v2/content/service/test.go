package service

import (
	"context"
	"fmt"
	"unsafe"

	"github.com/gomodule/redigo/redis"
	"github.com/liov/hoper/go/v2/initialize"
	"github.com/liov/hoper/go/v2/content/dao"
	model "github.com/liov/hoper/go/v2/protobuf/content"
	"github.com/liov/hoper/go/v2/protobuf/utils/request"
	"github.com/liov/hoper/go/v2/utils/encoding/json"
	"github.com/liov/hoper/go/v2/utils/net/http/tailmon"
	"github.com/liov/hoper/go/v2/utils/net/http/websocket"
)

type TestService struct {
	model.UnimplementedTestServiceServer
}

func (*TestService) GC(ctx context.Context, req *model.GCReq) (*request.Empty, error) {
	//address:= strconv.FormatUint()
	init := (*initialize.Init)(unsafe.Pointer(uintptr(req.Address)))
	fmt.Println(*init)
	return &request.Empty{}, nil
}

func (*TestService) Restart(ctx context.Context, req *request.Empty) (*request.Empty, error) {
	tailmon.ReStart()
	return &request.Empty{}, nil
}

func (*TestService) GetChat(ctx context.Context, req *request.Empty) ([]websocket.SendMessage, error) {
	conn := dao.Dao.Redis.Get()
	defer conn.Close()
	data, err := redis.ByteSlices(conn.Do("LRANGE", "Chat", 0, -1))
	if err != nil {
		return nil, err
	}
	var messages []websocket.SendMessage

	for _, v := range data {
		var message websocket.SendMessage
		json.Unmarshal(v, &message)
		messages = append(messages, message)
	}
	return messages, nil
}
