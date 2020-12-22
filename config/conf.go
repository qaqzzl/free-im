package config

import (
	"fmt"
	"free-im/pkg/logger"
	"go.uber.org/zap"
	"os"
	"github.com/spf13/viper"
)

var (
	CommonConf commonConf
	LogicConf logicConf
	ConnConf  connConf
	WSConf    wsConf
	HttpConf  httpConf
)
//
type commonConf struct {
	MySQL            string
	NSQIP            string
	RedisIP          string
}

// logic配置
type logicConf struct {
	MySQL            string
	NSQIP            string
	RedisIP          string
	RPCIntListenAddr string
	ConnRPCAddrs     string
}

// conn配置
type connConf struct {
	TCPListenAddr string
	RPCListenAddr string
	LogicRPCAddrs string
}

// WS配置
type wsConf struct {
	WSListenAddr  string
	RPCListenAddr string
	LogicRPCAddrs string
}

// Http配置
type httpConf struct {
	WSListenAddr  string
	RPCListenAddr string
	LogicRPCAddrs string
}

func init() {
	viper.SetConfigName("free")  // 配置文件名
	viper.SetConfigType("yaml") // 配置文件类型，可以是yaml、json、xml。。。
	viper.AddConfigPath(".")  // 配置文件路径
	err := viper.ReadInConfig()  // 读取配置文件信息
	if err != nil{
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	env := viper.Get("RunSetting")
	switch env {
	case "dev":
		logger.Leavel = zap.DebugLevel
		logger.Target = logger.File
	case "pre":
		logger.Leavel = zap.DebugLevel
		logger.Target = logger.File
	case "prod":
		logger.Leavel = zap.InfoLevel
		logger.Target = logger.File
	default:
		logger.Leavel = zap.DebugLevel
		logger.Target = logger.Console
	}

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
		LogicRPCAddrs: "127.0.0.1:50000",
	}

	WSConf = wsConf{
		WSListenAddr:  ":8081",
		RPCListenAddr: ":60001",
		LogicRPCAddrs: "127.0.0.1:50000",
	}
}
