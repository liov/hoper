package dingding

import "testing"

func TestRobot(t *testing.T) {
	SendRobotTextMessageWithSecret("xx", "xx", "hello world")
	SendRobotMarkDownMessageWithSecret("xxx", "xx", "xxx", "hello world")
}
