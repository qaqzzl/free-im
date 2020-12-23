package id

import (
	"github.com/gomodule/redigo/redis"
	"strconv"
)

type logger interface {
	Error(error)
}

// Logger Log接口，如果设置了Logger，就使用Logger打印日志，如果没有设置，就使用内置库log打印日志
var Logger logger

func Init(rdb *redis.Pool) {
	initChatroomID(rdb)
	initFreeID(rdb)
	initUID(rdb)
}

func byteUintToint64(bs []uint8) int64 {
	ba := []byte{}
	for _, b := range bs {
		ba = append(ba, byte(b))
	}
	in, _ := strconv.Atoi(string(ba))
	return int64(in)
}

func byteUintToString(bs []uint8) string {
	ba := []byte{}
	for _, b := range bs {
		ba = append(ba, byte(b))
	}
	return string(ba)
}