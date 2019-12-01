package model

import (
	"encoding/json"
	"net"
)
// 消息协议
const (
	MotionSignIn         			= 1 		// 设备登录
	MotionSignInACK      			= 2 		// 设备登录回执
	MotionSyncTrigger    			= 3 		// 消息同步触发
	MotionHeadbeat       			= 4 		// 心跳
	MotionHeadbeatACK    			= 5 		// 心跳回执
	MotionMessageSend    			= 6 		// 消息发送
	MotionMessageSendACK 			= 7 		// 消息发送回执
	MotionMessage        			= 8 		// 消息投递
	MotionMessageACK     			= 9 		// 消息投递回执
	MotionAuth		     			= 10		// 连接认证
	MotionQuit		     			= 11		// 客户端退出
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
	Code    int	`json:"code"`					// 消息类型
	MessageId string `json:"message_id"`		// 消息ID
	Content interface{} `json:"content"`		// 消息体
	Motion	int		`json:"motion"`				// 操作
}

//认证信息
type auth struct {
	DeviceID		string
	UserID			string
	AccessToken		string
	IsAuth			bool
}

type Context struct {
	ConnSocket		net.Conn				// 底层socket连接
	isClosed 		bool
	closeChan 		chan byte  				// 关闭通知

	Message			MessagePackage			//当前消息包

	InChan 			chan *MessagePackage	// 读队列 (入)
	OutChan 		chan *MessagePackage 	// 写队列 (出)
	Auth 			auth					// 认证信息
}


type Response struct {
	Code	int `json:"code"`
	Msg		string	`json:"msg"`
	Data	messagePackage `json:"packages"`
}

func (cxt *Context) Response (resp Response) {
	res,_ := json.Marshal(resp)
	cxt.ConnSocket.Write(res)
}