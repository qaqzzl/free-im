package tcp

import "net"

// IM连接信息
type Context struct {
	TcpConn net.Conn
	InChan  chan *[]byte	 // 入chan
	OutChan chan *[]byte	 // 出chan
	DeviceID string          // 设备id
	UserID   string          // 用户id
	IsAuth	bool			 // 是否登录(认证)
}
