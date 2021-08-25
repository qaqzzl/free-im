package main

import (
	api_tcp_conn "free-im/api/tcp_conn"
	"free-im/config"
	"free-im/internal/tcp_conn"
	"free-im/pkg/library/cache/redis"
	"free-im/pkg/rpc_client"
	"free-im/pkg/service"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	// net/http/pprof
	go func() {
		http.ListenAndServe("0.0.0.0:8899", nil)
	}()

	// 初始化公共服务
	service.Init(redis.NewPool(redis.Config{Dial: config.CommonConf.RedisIP, Auth: config.CommonConf.RedisAuth}))

	// 启动rpc服务
	go func() {
		api_tcp_conn.StartRPCServer()
	}()

	// 初始化 rpc 客户端
	rpc_client.InitLogic(config.TCPConnConf.LogicRPCAddrs)

	// 启动长链接服务器
	server := tcp_conn.NewTCPServer(config.TCPConnConf.TCPListenAddr, 10)
	server.Start()
}
