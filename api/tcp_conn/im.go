package tcp_conn

import (
	"github.com/orcaman/concurrent-map"
	"net"
)

//
//  自定义消息结构 , {消息操作类型(必选)}{消息体(非必选)}
//

const Version	int32 = 100000

// 消息操作类型 128
const (
	ActionMessageID         	int8	= 0 		// 消息ID
	ActionSignIn         		int8	= 1 		// 设备登录
	ActionSignInACK      		int8	= 2 		// 设备登录回执
	ActionSyncTrigger    		int8	= 3 		// 消息同步触发
	ActionMessageRead    		int8	= 4 		// 消息接收
	ActionMessageSend        	int8	= 5 		// 消息投递
	ActionMessageACK 			int8	= 6 		// 消息回执
	ActionClientACK     		int8	= 7 		// 客户端回执
	ActionAuth		     		int8	= 10		// 连接认证
	ActionQuit		     		int8	= 11		// 客户端退出
	ActionHeadbeat       		int8	= 99 		// 心跳
	ActionHeadbeatACK    		int8	= 100 		// 心跳回执
)

// 消息码(消息类型)
const (
	MessageCodeText         	= 1 		// 普通文本消息
	MessageCodePhiz     		= 2 		// 图片消息
	MessageCodeImage     		= 3 		// 位置
	MessageCodeVideo     		= 4 		// 视频消息
	MessageCodeFile     		= 5 		// 文件消息
)

// 会话类型(聊天室类型), 单聊、群聊、系统消息、聊天室、客服
const (
	ChatroomTypeSingle		int8	= 0			// 单聊
	ChatroomTypeGroup		int8	= 1			// 群聊
)

// Package 消息包
type MessagePackage struct {
	Code    	int	`json:"code"`					// 消息码(类型)
	ChatroomId	string `json:"chatroom_id"`			// 聊天室ID
	Content 	string `json:"content"`				// 消息体
	MessageId 	string `json:"message_id"`			// 消息ID
	UserId 		string `json:"user_id"`				// 用户ID
	ClientType 	string `json:"client_type"`			// 客户端类型
	MessageSendTime 	int64 `json:"message_send_time"`			// 消息发送时间, 服务器接收到的时间算
}

type Package struct {
	Version int32
	Action int8
	SequenceId int32
	BodyLength int32
	BodyData []byte
}


//认证信息
type AuthMessage struct {
	DeviceID		string `json:"device_id"`
	UserID			string `json:"user_id"`
	AccessToken		string `json:"access_token"`
	DeviceType 		string `json:"device_type"`		// 设备类型, 移动端:mobile , PC端:pc
	ClientType 		string `json:"client_type"`		// 客户端类型, android, ios,

}

// 连接用户客户端结构体
// DeviceType 设备类型, 移动端:mobile , PC端:pc
type ClientDevice struct {
	ClientType string		// 客户端类型, (android, ios) | (windows, mac, linux)
	DeviceID string
	Conn net.Conn
	Context *Context
}

//var SocketConnPool = make(map[string]map[string] ClientConn)		// 这是不支持并发的
//[user_id][DeviceType]ClientConn
var SocketConnPool = cmap.New()			//解决map并发读写

// 消息重发结构体

// 消息重发
// [user_id][DeviceType]