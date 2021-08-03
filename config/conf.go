package config

import (
	"flag"
	"fmt"
	"free-im/pkg/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	CommonConf  commonConf
	LogicConf   logicConf
	TCPConnConf tcpConnConf
	WSConnConf  wsConnConf
	HttpConf    httpConf
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

	QQAuthAppID string

	QiniuAccessKey string
	QiniuSecretKey string
}

// logic配置
type logicConf struct {
	RPCListenAddr string
	ConnRPCAddrs  string
}

// tcp conn配置
type tcpConnConf struct {
	TCPListenAddr string
	RPCListenAddr string
	LogicRPCAddrs string
}

// WS conn 配置
type wsConnConf struct {
	WSListenAddr  string
	RPCListenAddr string
	LogicRPCAddrs string
}

// Http配置
type httpConf struct {
	HttpListenAddr string
	RPCListenAddr  string
	LogicRPCAddrs  string
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
		// author local dev
		viper.AddConfigPath("/Users/zerozz/work/project/free-im/free-im/.") // 配置文件路径
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
		MySQL:     viper.GetString("MySQL"),
		NSQIP:     "",
		RedisIP:   viper.GetString("RedisIP"),
		RedisAuth: viper.GetString("RedisAuth"),

		AliYunSmsAccessKeyID:     viper.GetString("AliYunSmsAccessKeyID"),
		AliYunSmsAccessKeySecret: viper.GetString("AliYunSmsAccessKeySecret"),

		QQAuthAppID: viper.GetString("QQAuthAppID"),

		QiniuAccessKey: viper.GetString("QiniuAccessKey"),
		QiniuSecretKey: viper.GetString("QiniuSecretKey"),
	}

	LogicConf = logicConf{
		RPCListenAddr: viper.GetString("LogicConf.RPCListenAddr"),
		ConnRPCAddrs:  viper.GetString("LogicConf.ConnRPCAddrs"),
	}

	TCPConnConf = tcpConnConf{
		TCPListenAddr: viper.GetString("TCPConnConf.TCPListenAddr"),
		RPCListenAddr: viper.GetString("TCPConnConf.RPCListenAddr"),
		LogicRPCAddrs: viper.GetString("TCPConnConf.LogicRPCAddrs"),
	}

	WSConnConf = wsConnConf{
		WSListenAddr:  viper.GetString("WSConnConf.WSListenAddr"),
		RPCListenAddr: viper.GetString("WSConnConf.RPCListenAddr"),
		LogicRPCAddrs: viper.GetString("WSConnConf.LogicRPCAddrs"),
	}

	HttpConf = httpConf{
		HttpListenAddr: viper.GetString("HttpConf.HttpListenAddr"),
		RPCListenAddr:  viper.GetString("HttpConf.RPCListenAddr"),
		LogicRPCAddrs:  viper.GetString("HttpConf.LogicRPCAddrs"),
	}
}
