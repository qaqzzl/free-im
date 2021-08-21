package service

import (
	"free-im/pkg/service/id"
	"free-im/pkg/service/user"
	"github.com/gomodule/redigo/redis"
)

type logger interface {
	Error(error)
}

// Logger Log接口，如果设置了Logger，就使用Logger打印日志，如果没有设置，就使用内置库log打印日志
var Logger logger

func Init(rdb *redis.Pool) {
	id.InitChatroomID(rdb)
	id.InitFreeID(rdb)
	id.InitUID(rdb)
	user.InitUser(rdb)
}
