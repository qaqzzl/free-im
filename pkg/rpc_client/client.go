package rpc_client

import (
	"free-im/pkg/logger"
	"free-im/pkg/protos/pbs"
	"google.golang.org/grpc"
)

var (
	LogicInit   pbs.LogicInitClient
	ConnectInit pbs.ConnInitClient
)

func InitLogicInit(addr string) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		logger.Sugar.Error(err)
		panic(err)
	}
	defer conn.Close()

	LogicInit = pbs.NewLogicInitClient(conn)
}

func InitConnInit(addr string) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		logger.Sugar.Error(err)
		panic(err)
	}

	ConnectInit = pbs.NewConnInitClient(conn)
}
