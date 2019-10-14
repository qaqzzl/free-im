package db

import "net"

// 消息协议
const (
	CodeSignIn         = 1 // 设备登录
	CodeSignInACK      = 2 // 设备登录回执
	CodeSyncTrigger    = 3 // 消息同步触发
	CodeHeadbeat       = 4 // 心跳
	CodeHeadbeatACK    = 5 // 心跳回执
	CodeMessageSend    = 6 // 消息发送
	CodeMessageSendACK = 7 // 消息发送回执
	CodeMessage        = 8 // 消息投递
	CodeMessageACK     = 9 // 消息投递回执
)

const (
	MotionAuth			= 1		//连接认证
	MotionSendMessage	= 2		//客户端发送消息
	MotionPullMessage	= 3		//客户端拉取消息
	MotionQuit			= 4		//客户端退出
)

// Package 消息包
type MessagePackage struct {
	Code    int	`json:"code"`					//消息类型
	Content []byte `json:"content"`				// 消息体
}
type Message struct {
	Motion			int		`json:"motion"`
	AccessToken		string	`json:"access_token"`
	Package			MessagePackage	`json:"package"`
	DeviceID		int	`json:"device_id"`
	UserID			int	`json:"user_id"`
}

//链接信息
type Connect struct {
	Conn	*net.Conn
	IsSignIn bool   // 是否登录
	DeviceId int64  // 设备id
}

var Connects	[100000][] *Connect
