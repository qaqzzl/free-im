package config

import (
	"free-im/pkg/logger"
	"go.uber.org/zap"
)

func initLocalConf() {
	LogicConf = logicConf{
		MySQL:            "root:root@tcp(127.0.0.1:3306)/free_im?charset=utf8mb4",
		NSQIP:            "127.0.0.1:4150",
		RedisIP:          "127.0.0.1:6379",
		RPCIntListenAddr: ":50000",
		ConnRPCAddrs:     "127.0.0.1:60000",
	}

	ConnConf = connConf{
		TCPListenAddr: ":1208",
		RPCListenAddr: ":60000",
		LocalAddr:     "127.0.0.1:60000",
		LogicRPCAddrs: "127.0.0.1:50000",
	}

	WSConf = wsConf{
		WSListenAddr:  ":8081",
		RPCListenAddr: ":60001",
		LocalAddr:     "127.0.0.1:60001",
		LogicRPCAddrs: "127.0.0.1:50000",
	}

	logger.Leavel = zap.DebugLevel
	logger.Target = logger.Console
}
