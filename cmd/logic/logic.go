package main

import (
	"free-im/api/logic"
	"free-im/config"
	"free-im/pkg/rpc_client"
)

func main() {
	// 初始化 rpc 客户端
	rpc_client.InitConn(config.LogicConf.ConnRPCAddrs)

	logic.StartRPCServer()
}
