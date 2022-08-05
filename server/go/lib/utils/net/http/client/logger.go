package client

import "time"

type LogCallback func(url, method, auth, reqBody, respBytes string, status int, process time.Duration)

type Logger interface {
	SetPrefix(string)
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}
