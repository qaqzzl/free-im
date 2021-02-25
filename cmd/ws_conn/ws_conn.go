package main

import (
	api_ws_conn "free-im/api/ws_conn"
	"free-im/config"
	"free-im/internal/ws_conn"
	"free-im/pkg/rpc_client"
)

func main() {
	// 启动rpc服务
	go func() {
		api_ws_conn.StartRPCServer()
	}()

	// 初始化 rpc 客户端
	rpc_client.InitLogicInit(config.TCPConnConf.LogicRPCAddrs)

	// 启动长链接服务器
	wsServer := ws_conn.NewWebSocketServer(config.WSConnConf.WSListenAddr)
	wsServer.Start()
}
