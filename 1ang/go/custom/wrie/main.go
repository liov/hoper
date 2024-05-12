package main

import (
	"fmt"
	"os"
)

type Message string

type Greeter struct {
	Message Message // <- adding a Message field
}

func (g Greeter) Greet() Message {
	return g.Message
}

type Event struct {
	Greeter Greeter // <- adding a Greeter field
}

func (e Event) Start() {
	msg := e.Greeter.Greet()
	fmt.Println(msg)
}

func NewMessage(phrase string) Message {
	return Message(phrase)
}
func NewGreeter(m Message) Greeter {
	return Greeter{Message: m}
}

func main() {
	e, err := InitializeEvent("hello")
	if err != nil {
		fmt.Printf("failed to create event: %s\n", err)
		os.Exit(2)
	}
	e.Start()
}

func NewEvent(g Greeter) Event {
	return Event{Greeter: g}
}
