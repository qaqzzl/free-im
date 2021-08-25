package main

import (
	api_ws_conn "free-im/api/ws_conn"
	"free-im/config"
	"free-im/internal/ws_conn"
	"free-im/pkg/library/cache/redis"
	"free-im/pkg/rpc_client"
	"free-im/pkg/service"
)

func main() {
	// 初始化公共服务
	service.Init(redis.NewPool(redis.Config{Dial: config.CommonConf.RedisIP, Auth: config.CommonConf.RedisAuth}))

	// 启动rpc服务
	go func() {
		api_ws_conn.StartRPCServer()
	}()

	// 初始化 rpc 客户端
	rpc_client.InitLogic(config.WSConnConf.LogicRPCAddrs)

	// 启动长链接服务器
	wsServer := ws_conn.NewWebSocketServer(config.WSConnConf.WSListenAddr)
	wsServer.Start()
}
