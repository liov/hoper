package service

import (
	"context"
	"fmt"
	"unsafe"

	"github.com/liov/hoper/go/v2/protobuf/utils/empty"
	"github.com/liov/hoper/go/v2/tailmon/initialize"
	"github.com/liov/hoper/go/v2/content/dao"
	model "github.com/liov/hoper/go/v2/protobuf/content"
	redisi "github.com/liov/hoper/go/v2/utils/dao/redis"
	"github.com/liov/hoper/go/v2/utils/encoding/json"
	"github.com/liov/hoper/go/v2/tailmon"
	"github.com/liov/hoper/go/v2/utils/net/http/websocket"
)

type TestService struct {
	model.UnimplementedTestServiceServer
}

func (*TestService) GC(ctx context.Context, req *model.GCReq) (*empty.Empty, error) {
	//address:= strconv.FormatUint()
	init := (*initialize.Init)(unsafe.Pointer(uintptr(req.Address)))
	fmt.Println(*init)
	return &empty.Empty{}, nil
}

func (*TestService) Restart(ctx context.Context, req *empty.Empty) (*empty.Empty, error) {
	tailmon.ReStart()
	return &empty.Empty{}, nil
}

func (*TestService) GetChat(ctx context.Context, req *empty.Empty) ([]websocket.SendMessage, error) {
	conn := dao.Dao.Redis.Conn(ctx)
	defer conn.Close()
	data, err := redisi.ByteSlices(dao.Dao.Redis.Do(ctx,"LRANGE", "Chat", 0, -1).Result())
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
