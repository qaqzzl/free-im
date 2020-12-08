package id

import "github.com/gomodule/redigo/redis"

func Init(rdb *redis.Pool) {
	InitChatroomID(rdb)
}
