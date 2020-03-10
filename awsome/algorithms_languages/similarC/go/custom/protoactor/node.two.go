package main

import (
	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	messages "test/custom/protoactor/message"
)

type MyActor struct{}

func (*MyActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.Echo:
		context.Send(msg.Sender, &messages.Response{
			SomeValue: "result",
		})
	}
}

func main() {
	remote.Start("localhost:8091")

	// register a name for our local actor so that it can be spawned remotely
	remote.Register("hello", actor.PropsFromProducer(func() actor.Actor { return &MyActor{} }))
	console.ReadLine()
}
