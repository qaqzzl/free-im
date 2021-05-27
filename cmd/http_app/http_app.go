package main

import (
	"free-im/api/http_app"
	"free-im/config"
	"free-im/pkg/id"
	"free-im/pkg/library/cache/redis"
	"free-im/pkg/rpc_client"
)

func main() {
	// 初始化ID 生成器
	id.Init(redis.NewPool(redis.Config{Dial: config.CommonConf.RedisIP, Auth: config.CommonConf.RedisAuth}))

	// 初始化 rpc 客户端
	rpc_client.InitLogic(config.HttpConf.LogicRPCAddrs)

	http_app.StartHttpServer()
}
