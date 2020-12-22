package config

import (
	"flag"
	"fmt"
	"free-im/pkg/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	CommonConf commonConf
	LogicConf  logicConf
	ConnConf   connConf
	WSConf     wsConf
	HttpConf   httpConf
)

//
type commonConf struct {
	MySQL          string
	NSQIP          string
	RedisIP        string
	RedisAuth      string
	HttpListenAddr string

	AliYunSmsAccessKeyID     string
	AliYunSmsAccessKeySecret string
}

// logic配置
type logicConf struct {
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
	var confPath string
	flag.StringVar(&confPath, "c", ".", "set config file")
	flag.Parse()
	if confPath != "." {
		viper.SetConfigFile(confPath)
	} else {
		viper.SetConfigName("free")   // 配置文件名
		viper.SetConfigType("yaml")   // 配置文件类型，可以是yaml、json、xml。。。
		viper.AddConfigPath(confPath) // 配置文件路径
	}

	err := viper.ReadInConfig() // 读取配置文件信息
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	env := viper.GetString("RunSetting")
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

	CommonConf = commonConf{
		MySQL:          viper.GetString("MySQL"),
		NSQIP:          "",
		RedisIP:        viper.GetString("RedisIP"),
		RedisAuth:      viper.GetString("RedisAuth"),
		HttpListenAddr: viper.GetString("HttpListenAddr"),

		AliYunSmsAccessKeyID:     viper.GetString("AliYunSmsAccessKeyID"),
		AliYunSmsAccessKeySecret: viper.GetString("AliYunSmsAccessKeySecret"),
	}

	LogicConf = logicConf{
		RPCIntListenAddr: viper.GetString("LogicConf.RPCIntListenAddr"),
		ConnRPCAddrs:     viper.GetString("LogicConf.ConnRPCAddrs"),
	}

	ConnConf = connConf{
		TCPListenAddr: viper.GetString("ConnConf.TCPListenAddr"),
		RPCListenAddr: viper.GetString("ConnConf.RPCListenAddr"),
		LogicRPCAddrs: viper.GetString("ConnConf.LogicRPCAddrs"),
	}

	WSConf = wsConf{
		WSListenAddr:  viper.GetString("WSConf.WSListenAddr"),
		RPCListenAddr: viper.GetString("WSConf.RPCListenAddr"),
		LogicRPCAddrs: viper.GetString("WSConf.LogicRPCAddrs"),
	}
}
