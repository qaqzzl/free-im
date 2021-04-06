package main

import (
	api_tcp_conn "free-im/api/tcp_conn"
	"free-im/config"
	"free-im/internal/tcp_conn"
	"free-im/pkg/rpc_client"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	go func() {
		http.ListenAndServe("0.0.0.0:8899", nil)
	}()

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
