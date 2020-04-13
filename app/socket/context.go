package socket

import (
	"bufio"
	"net"
	"sync"
)

type sendMessage struct {
	Conn net.Conn
	ReceiveUserID string
	Message []byte
}

// IM连接信息
type Context struct {
	TcpConn net.Conn
	r    *bufio.Reader
	InChan  chan sendMessage	// 入chan
	OutChan chan sendMessage	// 出chan
	DeviceID string          	// 设备id
	UserID   string          	// 用户id
	DeviceType string			// 设备类型, 移动端:mobile , PC端:pc
	ClientType string			// 客户端类型, android, ios,
	IsAuth	bool			 	// 是否认证(登录)
	IsConnStatus bool			// 连接状态
	closeChan 			chan byte  						// 关闭通知
	mutex sync.Mutex									// 避免重复关闭管道
}
