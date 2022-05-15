package connect

import (
	"bufio"
	"errors"
	"free-im/pkg/logger"
	"free-im/pkg/protos/pbs"
	"io"
	"net"
	"time"
)

// tcpServer TCP服务器
type tcpServer struct {
	Address            string // 端口
	AcceptGoroutineNum int    // 接收建立连接的goroutine数量
}

// NewTCPServer 创建TCP服务器
func NewTCPServer(address string, acceptGoroutineNum int) *tcpServer {
	TCPServer = &tcpServer{
		Address:            address,
		AcceptGoroutineNum: acceptGoroutineNum,
	}
	return TCPServer
}

var TCPServer *tcpServer

// Start 启动服务器
func (t *tcpServer) Start() {
	logger.Sugar.Info("tcp server start", t.Address)

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
	select {}
}

// Accept 接收客户端的TCP长连接的建立
func (t *tcpServer) Accept(listener *net.TCPListener) {
	for {
		tcpConn, err := listener.AcceptTCP()
		if err != nil {
			print("Fail to connect, %s\n", err)
			break
		}
		err = tcpConn.SetKeepAlive(true)
		if err != nil {
			print("Fail to connect, %s\n", err)
			break
		}
		conn := NewConnContext(tcpConn)
		go DoConn(conn)
	}
}

func NewConnContext(c *net.TCPConn) *Conn {
	reader := bufio.NewReader(c)
	return &Conn{
		CoonType:      ConnTypeTCP,
		TCP:           c,
		TCPReader:     reader,
		TCPReadTicker: time.NewTicker(time.Millisecond * 100),
	}
}

// DoConn 处理TCP连接
func DoConn(conn *Conn) {
	conn.HandleConnect()
	for {
		if mp, err := read(conn); err != nil {
			logger.Sugar.Error(err)
			conn.Close()
			break
		} else {
			conn.HandlePackage(mp)
		}
	}
}

func read(conn *Conn) (mp pbs.MsgPackage, err error) {
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
			<-conn.TCPReadTicker.C
			continue
		}
		return mp, err
	}
}
