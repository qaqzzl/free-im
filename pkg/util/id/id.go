package id

import "github.com/gomodule/redigo/redis"

func Init(rdb redis.Conn) {
	InitChatroomID(rdb)
}
