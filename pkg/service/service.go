package service

import "github.com/gomodule/redigo/redis"

type service struct {
	redisPool *redis.Pool
}

var Service service

func initService(rdb *redis.Pool) error {
	Service = service{
		redisPool: rdb,
	}
	return nil
}
