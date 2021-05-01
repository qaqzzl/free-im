package dao

import (
	"free-im/config"
	"free-im/pkg/cache/redis"
	"free-im/pkg/db/orm"
	redis2 "github.com/gomodule/redigo/redis"
	"gorm.io/gorm"
)

type dao struct {
	db    *gorm.DB
	redis *redis2.Pool
}

var Dao = dao{
	redis: redis.NewPool(redis.Config{Dial: config.CommonConf.RedisIP, Auth: config.CommonConf.RedisAuth}),
	db:    orm.NewMySQL(&orm.Config{DSN: config.CommonConf.MySQL}),
}

func (d *dao) Ris() redis2.Conn {
	return Dao.redis.Get()
}

func (d *dao) DB() *gorm.DB {
	return d.db
}
