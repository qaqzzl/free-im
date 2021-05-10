package main

import (
	"free-im/api/http_app"
	"free-im/config"
	"free-im/pkg/id"
	"free-im/pkg/library/cache/redis"
)

func main() {
	// 初始化ID 生成器
	id.Init(redis.NewPool(redis.Config{Dial: config.CommonConf.RedisIP, Auth: config.CommonConf.RedisAuth}))

	http_app.StartHttpServer()
}
