package rpc_client

import (
	"fmt"
	"free-im/pkg/logger"
	"free-im/pkg/protos/pbs"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	LogicInit   pbs.LogicInitClient
	ConnectInit pbs.ConnInitClient
)

func InitLogic(addr string) {
	conn, err := grpc.DialContext(context.TODO(), addr, grpc.WithInsecure(), grpc.WithUnaryInterceptor(interceptor))
	if err != nil {
		logger.Sugar.Error(err)
		panic(err)
	}
	//conn, err := grpc.Dial("127.0.0.1:50000", grpc.WithInsecure())
	//if err != nil {
	//	panic(err)
	//}
	////defer conn.Close()
	LogicInit = pbs.NewLogicInitClient(conn)
}

func InitConn(addr string) {
	conn, err := grpc.DialContext(context.TODO(), addr, grpc.WithInsecure(), grpc.WithUnaryInterceptor(interceptor),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, "addr")))
	if err != nil {
		logger.Sugar.Error(err)
		panic(err)
	}

	ConnectInit = pbs.NewConnInitClient(conn)
}

func InitMonitor(addr string) {
	conn, err := grpc.DialContext(context.TODO(), addr, grpc.WithInsecure(), grpc.WithUnaryInterceptor(interceptor))
	if err != nil {
		logger.Sugar.Error(err)
		panic(err)
	}
	LogicInit = pbs.NewLogicInitClient(conn)
}
