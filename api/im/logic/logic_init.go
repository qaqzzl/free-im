package logic

import (
	"context"
	"free-im/configs"
	"free-im/internal/im/logic/service"
	"free-im/pkg/logger"
	"free-im/pkg/protos/pbs"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

type LogicInitServer struct{}

// MessageReceive 消息接收
func (s *LogicInitServer) MessageReceive(ctx context.Context, req *pbs.MessageReceiveReq) (*pbs.MessageReceiveResp, error) {
	return &pbs.MessageReceiveResp{}, service.MessageReceive(ctx, *req)
}

// MessageReceive 消息接收
func (s *LogicInitServer) MessageACK(ctx context.Context, req *pbs.MessageACKReq) (*pbs.MessageACKResp, error) {
	return &pbs.MessageACKResp{}, service.MessageACK(ctx, *req)
}

// UnaryServerInterceptor 服务器端的单向调用的拦截器
func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	resp, err := handler(ctx, req)
	logger.Logger.Debug("interceptor", zap.Any("info", info), zap.Any("req", req), zap.Any("resp", resp))
	return resp, err
}

// StartRPCServer 启动rpc服务器
func StartRPCServer() {
	listener, err := net.Listen("tcp", config.LogicConf.RPCIntListenAddr)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer(grpc.UnaryInterceptor(UnaryServerInterceptor))
	pbs.RegisterLogicInitServer(server, &LogicInitServer{})
	logger.Logger.Debug("rpc服务已经开启")
	err = server.Serve(listener)
	if err != nil {
		panic(err)
	}
}
