package tcp_conn

import (
	"bufio"
	"context"
	"errors"
	"free-im/pkg/logger"
	"free-im/pkg/protos/pbs"
	"github.com/orcaman/concurrent-map"
	"io"
	"net"
	"time"
)

//var SocketConnPool = make(map[string]map[string] Context)		// 这是不支持并发的
//[user_id][DeviceType]Context
var SocketConnPool = cmap.New() //解决map并发读写

// 连接用户客户端结构体
// DeviceType 设备类型, 移动端:mobile , PC端:pc
type ClientDevice struct {
	ClientType string // 客户端类型, (android, ios) | (windows, mac, linux)
	DeviceID   string
	Context    *Context
}

// IM连接信息
type Context struct {
	TcpConn      net.Conn
	r            *bufio.Reader
	Version      int32
	WriteChan    chan sendMessage // 出chan
	ReadChan     chan sendMessage // 入chan
	DeviceID     string           // 设备id 简写 DID
	UserID       string           // 用户id 简写 UID
	DeviceType   string           // 设备类型, 移动端:mobile , PC端:pc
	ClientType   string           // 客户端类型, android, ios,
	IsAuth       bool             // 是否认证(登录)
	IsConnStatus bool             // 连接状态
}

type sendMessage struct {
	Conn    net.Conn
	Package pbs.MessagePackage
}

func NewConnContext(conn *net.TCPConn) *Context {
	reader := bufio.NewReader(conn)
	return &Context{
		TcpConn:   conn,
		r:         reader,
		WriteChan: make(chan sendMessage, 1000),
		ReadChan:  make(chan sendMessage, 1000),
	}
}

var readTicker = time.NewTicker(time.Millisecond * 100)

// DoConn 处理TCP连接
func (ctx *Context) DoConn() {
	ctx.HandleConnect()
	for {
		if mp, err := ctx.Read(); err != nil {
			logger.Sugar.Error(err)
			ctx.TcpConn.Close()
			break
		} else {
			ctx.HandlePackage(mp)
		}
	}
}

// HandleConnect 建立连接
func (ctx *Context) HandleConnect() {
	ctx.IsConnStatus = true
}

func (ctx *Context) Read() (mp pbs.MessagePackage, err error) {
	var waitingReadCount = 0
	for {
		mp, err := Protocol.Decode(ctx)
		if err == io.EOF || err != nil {
			return mp, err
		}

		if waitingReadCount > 300 {
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

func (c *Context) Write(conn net.Conn, mp pbs.MessagePackage) (int, error) {
	if b, err := Protocol.Encode(mp); err != nil {
		return 0, err
	} else {
		nn, err := conn.Write(b)
		if err != nil {
			return 0, err
		}
		return nn, nil
	}
}

// HandlePackage 处理消息包
func (ctx *Context) HandlePackage(mp pbs.MessagePackage) {
	Handler.Handler(ctx, mp)
}

//send message handle
func (ctx *Context) SendMessage(conn net.Conn, mp pbs.MessagePackage) (n int, err error) {
	if n, err = ctx.Write(conn, mp); err == nil {
		return n, nil
	} else {
		return n, err
	}
}

// 投递消息
func DeliverMessageByUID(ctx context.Context, req *pbs.DeliverMessageReq) error {
	// 获取设备对应的TCP连接
	ctxconns := loadConnsByUID(req.UserId)
	if ctxconns == nil {
		logger.Sugar.Warn("ctx id nil")
		return nil
	}
	for _, ctxconn := range ctxconns {
		// 发送消息
		ctxconn.SendMessage(ctxconn.TcpConn, *req.Message)
	}
	return nil
}

// 投递消息
func DeliverMessageByUIDAndDID(ctx context.Context, req *pbs.DeliverMessageReq) error {
	// 获取设备对应的TCP连接
	ctxconns := loadConnsByUID(req.UserId)
	if ctxconns == nil {
		logger.Sugar.Warn("ctx id nil")
		return nil
	}
	for _, ctxconn := range ctxconns {
		if req.DeviceId == ctxconn.DeviceID {
			// 发送消息
			ctxconn.SendMessage(ctxconn.TcpConn, *req.Message)
		}
	}
	return nil
}

func DeliverMessageByUIDAndNotDID(ctx context.Context, req *pbs.DeliverMessageReq) error {
	// 获取设备对应的TCP连接
	ctxconns := loadConnsByUID(req.UserId)
	if ctxconns == nil {
		logger.Sugar.Warn("ctx id nil")
		return nil
	}
	for _, ctxconn := range ctxconns {
		if ctxconn.DeviceID == req.DeviceId {
			continue
		}
		// 发送消息
		ctxconn.SendMessage(ctxconn.TcpConn, *req.Message)
	}
	return nil
}

// 关闭链接
func Close(ctx *Context) {
	if ctx.IsConnStatus != false {
		if user_map, ok := SocketConnPool.Get(ctx.UserID); ok {
			user_map.(cmap.ConcurrentMap).Remove(ctx.DeviceType)
			SocketConnPool.Set(ctx.UserID, user_map)
		}
	}
	ctx.IsConnStatus = false
	ctx.TcpConn.Close()
}

// store 存储
func storeConn(userId string, ctx *Context) {
}

// load 获取链接
func loadConnsByUID(UserID string) (ctxs []*Context) {
	tmp, ok := SocketConnPool.Get(UserID)
	if ok && tmp.(cmap.ConcurrentMap).Count() > 0 {
		for _, vo := range tmp.(cmap.ConcurrentMap).Items() {
			device := vo.(ClientDevice)
			ctxs = append(ctxs, device.Context)
		}
	}
	return ctxs
}

// delete 删除
func deleteConn(userId string) {
}

func loadConnsByUIDAndDID(UserID string, DeviceID string) (ctx *Context) {
	tmp, ok := SocketConnPool.Get(UserID)
	if ok && tmp.(cmap.ConcurrentMap).Count() > 0 {
		for _, vo := range tmp.(cmap.ConcurrentMap).Items() {
			device := vo.(ClientDevice)
			if device.DeviceID == DeviceID {
				ctx = device.Context
				break
			}
		}
	}
	return ctx
}
