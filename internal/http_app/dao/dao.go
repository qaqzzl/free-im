package dao

import (
	"free-im/config"
	"free-im/pkg/library/cache/redis"
	orm2 "free-im/pkg/library/database/orm"
	"free-im/pkg/logger"
	redis2 "github.com/gomodule/redigo/redis"
	"gorm.io/gorm"
)

type dao struct {
	db    *gorm.DB
	redis *redis2.Pool
}

var Dao = dao{
	redis: redis.NewPool(redis.Config{Dial: config.CommonConf.RedisIP, Auth: config.CommonConf.RedisAuth}),
	db:    orm2.NewMySQL(&orm2.Config{DSN: config.CommonConf.MySQL}),
}

var redis_conn redis2.Conn

func (d *dao) Ris() redis2.Conn {
	var err error
	if redis_conn != nil {
		_, err = redis_conn.Do("PING")
	}
	if redis_conn == nil || err != nil {
		logger.Sugar.Info("redis 获取新链接")
		redis_conn = Dao.redis.Get()
	}
	return redis_conn
}

func (d *dao) DB() *gorm.DB {
	return d.db
}

func Paginate(page int, prepage int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case prepage > 100:
			prepage = 100
		case prepage <= 0:
			prepage = 10
		}

		offset := (page - 1) * prepage
		return db.Offset(offset).Limit(prepage)
	}
}
