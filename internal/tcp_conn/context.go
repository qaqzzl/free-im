package tcp_conn

import (
	"bufio"
	"net"
)

type sendMessage struct {
	Conn net.Conn
	Package Package
}

// IM连接信息
type Context struct {
	TcpConn net.Conn
	r    *bufio.Reader
	WriteChan chan sendMessage	// 出chan
	ReadChan  chan sendMessage	// 入chan
	DeviceID string          	// 设备id
	UserID   string          	// 用户id
	DeviceType string			// 设备类型, 移动端:mobile , PC端:pc
	ClientType string			// 客户端类型, android, ios,
	IsAuth	bool			 	// 是否认证(登录)
	IsConnStatus bool			// 连接状态
	Status bool
}
