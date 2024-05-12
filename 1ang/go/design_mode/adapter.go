package main

import "fmt"

type Player interface {
	play(string)
}

type VLC struct{}

func (v *VLC) play(name string) {
	fmt.Println("播放VLC文件: " + name)
}

type MP4 struct{}

func (m *MP4) play(name string) {
	fmt.Println("播放MP4文件: " + name)
}

type Adapter struct {
	Player
}

func (a *Adapter) get(typ string) {
	if typ == "vlc" {
		a.Player = &VLC{}
	} else {
		a.Player = &MP4{}
	}
}

func (a *Adapter) play(name string, typ string) {
	a.get(typ)
	a.Player.play(name)
}

func main() {
	var player Adapter
	player.play("mp4", "alone.mp4")
	player.play("vlc", "far far away.vlc")
}
