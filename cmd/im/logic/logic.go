package main

import (
	"free-im/api/im/logic"
	"free-im/configs"
	"free-im/internal/im/logic/dao"
	"free-im/pkg/rpc_client"
)

func main() {
	// 初始化 redis
	dao.InitRedis()

	// 初始化 rpc 客户端
	rpc_client.InitConnInit(config.LogicConf.ConnRPCAddrs)

	logic.StartRPCServer()
}
