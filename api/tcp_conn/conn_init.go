package tcp_conn

import (
	"context"
	"free-im/config"
	"free-im/internal/tcp_conn"
	"free-im/pkg/logger"
	"free-im/pkg/protos/pbs"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

type ConnInitServer struct{}

// Message 投递消息
func (s *ConnInitServer) DeliverMessageByUID(ctx context.Context, req *pbs.DeliverMessageReq) (*pbs.DeliverMessageResp, error) {
	return &pbs.DeliverMessageResp{}, tcp_conn.Handler.DeliverMessageByUID(req.UserId, *req.Message)
}

// Message 投递消息
func (s *ConnInitServer) DeliverMessageByUIDAndDID(ctx context.Context, req *pbs.DeliverMessageReq) (*pbs.DeliverMessageResp, error) {
	return &pbs.DeliverMessageResp{}, tcp_conn.Handler.DeliverMessageByUIDAndDID(req.UserId, req.DeviceId, *req.Message)
}

// Message 投递消息
func (s *ConnInitServer) DeliverMessageByUIDAndNotDID(ctx context.Context, req *pbs.DeliverMessageReq) (*pbs.DeliverMessageResp, error) {
	return &pbs.DeliverMessageResp{}, tcp_conn.Handler.DeliverMessageByUIDAndNotDID(req.UserId, req.DeviceId, *req.Message)
}

// UnaryServerInterceptor 服务器端的单向调用的拦截器
func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	resp, err := handler(ctx, req)
	logger.Logger.Debug("interceptor", zap.Any("info", info), zap.Any("req", req), zap.Any("resp", resp))
	return resp, err
}

// StartRPCServer 启动rpc服务器
func StartRPCServer() {
	listener, err := net.Listen("tcp", config.TCPConnConf.RPCListenAddr)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer(grpc.UnaryInterceptor(UnaryServerInterceptor))
	pbs.RegisterConnInitServer(server, &ConnInitServer{})
	logger.Logger.Debug("rpc服务已经开启")
	err = server.Serve(listener)
	if err != nil {
		panic(err)
	}
}
