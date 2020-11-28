package config

import (
	"free-im/pkg/logger"

	"go.uber.org/zap"
)

func initPreConf() {
	LogicConf = logicConf{
		MySQL:            "root:liu123456@tcp(localhost:3306)/gim?charset=utf8&parseTime=true",
		NSQIP:            "127.0.0.1:4150",
		RedisIP:          "127.0.0.1:6379",
		RPCIntListenAddr: ":50000",
		ConnRPCAddrs:     "addrs:///127.0.0.1:60000,127.0.0.1:60001",
	}

	ConnConf = connConf{
		TCPListenAddr: ":1208",
		RPCListenAddr: ":60000",
		LocalAddr:     "127.0.0.1:60000",
		LogicRPCAddrs: "addrs:///127.0.0.1:50000",
	}

	WSConf = wsConf{
		WSListenAddr:  ":8081",
		RPCListenAddr: ":60001",
		LocalAddr:     "127.0.0.1:60001",
		LogicRPCAddrs: "addrs:///127.0.0.1:50000",
	}

	logger.Leavel = zap.DebugLevel
	logger.Target = logger.File
}
