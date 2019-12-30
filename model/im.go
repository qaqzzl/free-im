package model

import (
	"net"
)

// 消息协议
const (
	ActionSignIn         			= 1 		// 设备登录
	ActionSignInACK      			= 2 		// 设备登录回执
	ActionSyncTrigger    			= 3 		// 消息同步触发
	ActionHeadbeat       			= 4 		// 心跳
	ActionHeadbeatACK    			= 5 		// 心跳回执
	ActionMessageSend    			= 6 		// 消息发送
	ActionMessageSendACK 			= 7 		// 消息发送回执
	ActionMessage        			= 8 		// 消息投递
	ActionMessageACK     			= 9 		// 消息投递回执
	ActionAuth		     			= 10		// 连接认证
	ActionQuit		     			= 11		// 客户端退出
)

// 消息码(消息类型)
const (
	CodeText         	= 1 		// 普通文本消息
	CodePhiz     		= 2 		// 表情消息
	CodeImage     		= 3 		// 图片消息
	CodeVideo     		= 4 		// 视频消息
	CodeFile     		= 5 		// 文件消息
)

// Package 消息包
type MessagePackage struct {
	Code    	int	`json:"code"`					// 消息码(类型)
	ChatroomId	string `json:"chatroom_id"`			// 聊天室ID
	Content 	[]byte `json:"content"`				// 消息体
	Action		int		`json:"action"`				// 操作
}

//认证信息
type AuthMessage struct {
	DeviceID		string `json:"device_id"`
	UserID			string `json:"user_id"`
	AccessToken		string `json:"access_token"`
}
							//[user_id][device_id]
var SocketConnPool = make(map[string]map[string] net.Conn)