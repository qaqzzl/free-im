package main

import (
	"free-im/api/http_app"
	"free-im/config"
	"free-im/pkg/library/cache/redis"
	"free-im/pkg/rpc_client"
	"free-im/pkg/service"
)

func main() {
	// 初始化公共服务
	service.Init(redis.NewPool(redis.Config{Dial: config.CommonConf.RedisIP, Auth: config.CommonConf.RedisAuth}))

	// 初始化 rpc 客户端
	rpc_client.InitLogic(config.HttpConf.LogicRPCAddrs)

	http_app.StartHttpServer()
}
