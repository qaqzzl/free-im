package tcp_conn

import (
	"free-im/pkg/logger"
	cmap "github.com/orcaman/concurrent-map"
	"net"
)

// sync.Map
// [user_id][DeviceType]Context
var ServerConnPool = cmap.New() //解决map并发读写

// TCPServer TCP服务器
type TCPServer struct {
	Address            string // 端口
	AcceptGoroutineNum int    // 接收建立连接的goroutine数量
}

// NewTCPServer 创建TCP服务器
func NewTCPServer(address string, acceptGoroutineNum int) *TCPServer {
	return &TCPServer{
		Address:            address,
		AcceptGoroutineNum: acceptGoroutineNum,
	}
}

// Start 启动服务器
func (t *TCPServer) Start() {
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
func (t *TCPServer) Accept(listener *net.TCPListener) {
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
