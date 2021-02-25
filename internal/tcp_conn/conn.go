package tcp_conn

import (
	"bufio"
	"errors"
	"free-im/pkg/logger"
	"free-im/pkg/protos/pbs"
	"github.com/orcaman/concurrent-map"
	"io"
	"net"
	"time"
)

// IM连接信息
type Context struct {
	TcpConn    net.Conn
	r          *bufio.Reader
	Version    int32
	DeviceID   string // 设备id 简写 DID
	UserID     string // 用户id 简写 UID
	DeviceType string // 设备类型, 移动端:mobile , PC端:pc
	ClientType string // 客户端类型, (android, ios) | (windows, mac, linux)
	IsAuth     bool   // 是否认证(登录)
}

func NewConnContext(conn *net.TCPConn) *Context {
	reader := bufio.NewReader(conn)
	return &Context{
		TcpConn: conn,
		r:       reader,
	}
}

// DoConn 处理TCP连接
func (ctx *Context) DoConn() {
	ctx.HandleConnect()
	for {
		if mp, err := ctx.Read(); err != nil {
			logger.Sugar.Error(err)
			ctx.Close()
			break
		} else {
			ctx.HandlePackage(mp)
		}
	}
}

// HandleConnect 建立连接
func (ctx *Context) HandleConnect() {
	logger.Logger.Info("connect")
}

func (ctx *Context) Read() (mp pbs.MsgPackage, err error) {
	var readTicker = time.NewTicker(time.Millisecond * 100)
	var waitingReadCount = 0
	for {
		mp, err := Protocol.Decode(ctx)
		if err == io.EOF || err != nil {
			return mp, err
		}

		if waitingReadCount > 100 { // 10s
			return mp, errors.New("time out")
		}
		if mp.BodyData == nil {
			waitingReadCount++
			<-readTicker.C
			continue
		}
		return mp, err
	}
}

func (ctx *Context) Write(mp pbs.MsgPackage) (int, error) {
	if b, err := Protocol.Encode(mp); err != nil {
		return 0, err
	} else {
		nn, err := ctx.TcpConn.Write(b)
		if err != nil {
			return 0, err
		}
		return nn, nil
	}
}

func (ctx *Context) Close() {
	if user_map, ok := TCPServer.ServerConnPool.Get(ctx.UserID); ok {
		user_map.(cmap.ConcurrentMap).Remove(ctx.DeviceType)
		// ServerConnPool.Set(ctx.UserID, user_map)
	}
	ctx.TcpConn.Close()
}

// HandlePackage 处理消息包
func (ctx *Context) HandlePackage(mp pbs.MsgPackage) {
	Handler.Handler(ctx, mp)
}
