package dao

import (
	"database/sql"
	"free-im/pkg/cache/redis"
)

type dao struct {
	db    *sql.DB
	redis *redis.Pool
}

var Dao = dao{
	redis: redis.NewPool(),
}
