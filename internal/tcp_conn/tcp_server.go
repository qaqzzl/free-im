package tcp_conn

import (
	"free-im/pkg/logger"
	"free-im/pkg/service/user"
	cmap "github.com/orcaman/concurrent-map"
	"net"
	"strconv"
)

// tcpServer TCP服务器
type tcpServer struct {
	Address            string // 端口
	AcceptGoroutineNum int    // 接收建立连接的goroutine数量
	// sync.Map
	// [user_id][DeviceType]Context
	ServerConnPool cmap.ConcurrentMap // 链接池
}

// NewTCPServer 创建TCP服务器
func NewTCPServer(address string, acceptGoroutineNum int) *tcpServer {
	TCPServer = &tcpServer{
		Address:            address,
		AcceptGoroutineNum: acceptGoroutineNum,
		ServerConnPool:     cmap.New(),
	}
	return TCPServer
}

var TCPServer *tcpServer

// Start 启动服务器
func (t *tcpServer) Start() {
	addr, err := net.ResolveTCPAddr("tcp", t.Address)
	if err != nil {
		panic(err)
	}
	server, err := net.ListenTCP("tcp", addr)
	if err != nil {
		print("Fail to start server, %s\n", err)
	}
	for i := 0; i < t.AcceptGoroutineNum; i++ {
		go t.Accept(server)
	}
	logger.Sugar.Info("tcp server start")
	select {}
}

// Accept 接收客户端的TCP长连接的建立
func (t *tcpServer) Accept(listener *net.TCPListener) {
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			print("Fail to connect, %s\n", err)
			break
		}
		err = conn.SetKeepAlive(true)
		if err != nil {
			print("Fail to connect, %s\n", err)
			break
		}
		connContext := NewConnContext(conn)
		go connContext.DoConn()
	}
}

// load 获取链接
func (t *tcpServer) LoadConn(UserID string, DeviceID string) (conn *Conn) {
	tmp, ok := t.ServerConnPool.Get(UserID)
	if ok && tmp.(cmap.ConcurrentMap).Count() > 0 {
		for _, vo := range tmp.(cmap.ConcurrentMap).Items() {
			conn := vo.(*Conn)
			if conn.DeviceID == DeviceID {
				break
			}
		}
	}
	return conn
}

func (t *tcpServer) LoadConnsByUID(UserID int64) (conns []*Conn) {
	tmp, ok := t.ServerConnPool.Get(strconv.Itoa(int(UserID)))
	if ok && tmp.(cmap.ConcurrentMap).Count() > 0 {
		for _, vo := range tmp.(cmap.ConcurrentMap).Items() {
			conn := vo.(*Conn)
			conns = append(conns, conn)
		}
	}
	return conns
}

// store 存储
func (t *tcpServer) StoreConn(conn *Conn) {
	key := strconv.Itoa(int(conn.UserID))
	if tmp, ok := t.ServerConnPool.Get(key); ok {
		device_map := tmp.(cmap.ConcurrentMap)
		device_map.Set(conn.DeviceType, conn)
		t.ServerConnPool.Set(key, device_map)
	} else {
		device_map := cmap.New()
		device_map.Set(conn.DeviceType, conn)
		t.ServerConnPool.Set(key, device_map)
	}
	// 用户在线状态
	user.User.SetUserOnline(conn.UserID, true, conn.DeviceType)
}
