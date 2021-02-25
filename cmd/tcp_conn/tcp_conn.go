package main

import (
	api_tcp_conn "free-im/api/tcp_conn"
	"free-im/config"
	"free-im/internal/tcp_conn"
	"free-im/pkg/rpc_client"
)

func main() {
	// 启动rpc服务
	go func() {
		api_tcp_conn.StartRPCServer()
	}()

	// 初始化 rpc 客户端
	rpc_client.InitLogicInit(config.TcpConnConf.LogicRPCAddrs)

	// 启动长链接服务器
	server := tcp_conn.NewTCPServer(config.TcpConnConf.TCPListenAddr, 10)
	server.Start()
}
