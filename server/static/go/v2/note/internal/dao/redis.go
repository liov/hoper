package dao

import (
	"github.com/gomodule/redigo/redis"
)

type NoteRedis struct {
	redis.Conn
}

func (conn *NoteRedis) Close() {
	conn.Conn.Close()
}

func NewNoteRedis() *NoteRedis {
	conn := Dao.Redis.Get()
	return &NoteRedis{conn}
}
