package dao

import (
	"free-im/pkg/cache/redis"
	"free-im/pkg/db"
)

type dao struct {
	db    *db.Db
	redis *redis.Pool
}

var Dao = dao{
	redis: redis.NewPool(),
}
