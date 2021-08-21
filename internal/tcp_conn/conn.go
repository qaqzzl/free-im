package tcp_conn

import (
	"bufio"
	"errors"
	"free-im/pkg/logger"
	"free-im/pkg/protos/pbs"
	"free-im/pkg/service/user"
	cmap "github.com/orcaman/concurrent-map"
	"io"
	"net"
	"strconv"
	"time"
)

// IM连接信息
type Conn struct {
	c          net.Conn
	r          *bufio.Reader
	Version    int32
	DeviceID   string // 设备id 简写 DID
	UserID     int64  // 用户id 简写 UID
	DeviceType string // 设备类型, 移动端:mobile , PC端:pc
	ClientType string // 客户端类型, (android, ios) | (windows, mac, linux)
	IsAuth     bool   // 是否认证(登录)
	readTicker *time.Ticker
}

func NewConnContext(c *net.TCPConn) *Conn {
	reader := bufio.NewReader(c)
	return &Conn{
		c:          c,
		r:          reader,
		readTicker: time.NewTicker(time.Millisecond * 100),
	}
}

// DoConn 处理TCP连接
func (conn *Conn) DoConn() {
	conn.HandleConnect()
	for {
		if mp, err := conn.Read(); err != nil {
			logger.Sugar.Error(err)
			conn.Close()
			break
		} else {
			conn.HandlePackage(mp)
		}
	}
}

// HandleConnect 建立连接
func (conn *Conn) HandleConnect() {
	logger.Logger.Info("connect")
}

func (conn *Conn) Read() (mp pbs.MsgPackage, err error) {
	var waitingReadCount = 0
	for {
		mp, err := Protocol.Decode(conn)
		if err == io.EOF || err != nil {
			return mp, err
		}

		if waitingReadCount > 100 { // 10s
			return mp, errors.New("time out")
		}
		if mp.BodyData == nil {
			waitingReadCount++
			<-conn.readTicker.C
			continue
		}
		return mp, err
	}
}

func (conn *Conn) Write(mp pbs.MsgPackage) (int, error) {
	//logger.Sugar.Debug(mp)
	if b, err := Protocol.Encode(mp); err != nil {
		return 0, err
	} else {
		nn, err := conn.c.Write(b)
		if err != nil {
			return 0, err
		}
		return nn, nil
	}
}

func (conn *Conn) Close() {
	if conn.IsAuth {
		key := strconv.Itoa(int(conn.UserID))
		if user_map, ok := TCPServer.ServerConnPool.Get(key); ok {
			user_map.(cmap.ConcurrentMap).Remove(conn.DeviceType)
			TCPServer.ServerConnPool.Set(key, user_map)
		}
		// 用户在线状态
		user.User.SetUserOnline(conn.UserID, false, conn.DeviceType)
	}
	conn.readTicker.Stop()
	conn.c.Close()
}

// HandlePackage 处理消息包
func (conn *Conn) HandlePackage(mp pbs.MsgPackage) {
	Handler.Handler(conn, mp)
}
