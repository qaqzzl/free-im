package redis

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

// Config mysql config.
type Config struct {
	Dial string
	Auth string
}

func NewPool(cnf Config) *redis.Pool {
	return &redis.Pool{
		// 最大的激活连接数，表示同时最多有N个连接 ，为0事表示没有限制
		MaxActive: 0,
		//最大的空闲连接数，表示即使没有redis连接时依然可以保持N个空闲的连接，而不被清除，随时处于待命状态
		MaxIdle: 0,
		//最大的空闲连接等待时间，超过此时间后，空闲连接将被关闭
		IdleTimeout: 240 * time.Second,
		// 当链接数达到最大后是否阻塞，如果不的话，达到最大后返回错误
		Wait: false,
		//Dial 是创建链接的方法
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", cnf.Dial)
			if err != nil {
				return nil, err
			}
			if cnf.Auth != "" {
				if _, err := c.Do("AUTH", cnf.Auth); err != nil {
					c.Close()
					return nil, err
				}
			}

			return c, err
		},
		//TestOnBorrow 是一个测试链接可用性的方法
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}
