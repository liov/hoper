package service

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/liov/hoper/server/go/protobuf/message"
	"google.golang.org/protobuf/proto"
)

var messageService = &MessageService{}

type MessageService struct {
	message.UnimplementedMessageServer
}

func (*MessageService) Send(ctx context.Context, req *message.MQMessage) (*empty.Empty, error) {

	return &empty.Empty{}, nil
}

func (s *MessageService) Receive(ctx context.Context, req *message.ClientMessage) (*empty.Empty, error) {
	var err error
	switch req.Command {
	case message.ClientCmdJoinGroup:
		var req1 message.JoinGroupReq
		err = proto.Unmarshal(req.Payload, &req1)
		if err != nil {
			return nil, err
		}
		_, err = s.JoinGroup(ctx, &req1)
		if err != nil {
			return nil, err
		}
	}

	return &empty.Empty{}, nil
}

func (*MessageService) JoinGroup(ctx context.Context, req *message.JoinGroupReq) (*empty.Empty, error) {

	return &empty.Empty{}, nil
}
