package model

import (
	"github.com/orcaman/concurrent-map"
	"net"
)

//
//  自定义消息结构 , {消息操作类型(必选)}{消息体(非必选)}
//

// 消息操作类型
const (
	ActionSignIn         			= "01" 		// 设备登录
	ActionSignInACK      			= "02" 		// 设备登录回执
	ActionSyncTrigger    			= "03" 		// 消息同步触发
	ActionMessage    				= "04" 		// 消息
	ActionMessageACK 				= "05" 		// 消息回执
	ActionMessageSend        		= "08" 		// 消息投递
	ActionMessageSendACK     		= "09" 		// 消息投递回执
	ActionAuth		     			= "10"		// 连接认证
	ActionQuit		     			= "11"		// 客户端退出
	ActionHeadbeat       			= "99" 		// 心跳
	ActionHeadbeatACK    			= "" 		// 心跳回执
)

// 消息码(消息类型)
const (
	CodeText         	= 1 		// 普通文本消息
	CodePhiz     		= 2 		// 图片消息
	CodeImage     		= 3 		// 位置
	CodeVideo     		= 4 		// 视频消息
	CodeFile     		= 5 		// 文件消息
)

// Package 消息包
type MessagePackage struct {
	Code    	int	`json:"code"`						// 消息码(类型)
	ChatroomId	string `json:"chatroom_id"`			// 聊天室ID
	Content 	string `json:"content"`				// 消息体
	MessageId 	string `json:"message_id"`			// 客户端消息ID
	UserId 	string `json:"user_id"`					// 用户ID
}

// 服务端发送的 Package 消息包
type ServerSendMessagePackage struct {
	Code    	int	`json:"code"`										// 消息码(类型)
	ChatroomId	string `json:"chatroom_id"`							// 聊天室ID
	Content 	string `json:"content"`								// 消息体
	ClientMessageId 	string `json:"client_message_id"`			// 客户端消息ID
	ServerMessageId 	string `json:"server_message_id"`			// 服务端消息ID
	UserId 	string `json:"user_id"`									// 发送消息的用户ID
	MessageSendTime 	int64 `json:"message_send_time"`			// 消息发送时间, 服务器接收到的时间算
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
	ClientType string		// 客户端类型, android, ios,
	DeviceID string
	Conn net.Conn
}
							//[user_id][DeviceType]ClientConn
//var SocketConnPool = make(map[string]map[string] ClientConn)		// 这是不支持并发的
var SocketConnPool = cmap.New()			//解决map并发读写