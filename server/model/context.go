package model

import (
	"encoding/json"
	"net"
)

// Package 消息包
type messagePackage struct {
	code    int	`json:"code"`					// 消息类型
	content []byte `json:"content"`				// 消息体
}

//消息体
type Message struct {
	Motion			int		`json:"motion"`				// 操作
	packages		messagePackage	`json:"packages"`
}

//认证信息
type auth struct {
	DeviceID		int	`json:"device_id"`
	UserID			int	`json:"user_id"`
	AccessToken		string	`json:"access_token"`
}

type Context struct {
	ConnSocket		net.Conn				// 底层socket连接
	isClosed 		bool
	closeChan 		chan byte  				// 关闭通知

	inChan 			chan *Message			// 读队列
	outChan 		chan *messagePackage 	// 写队列
	auth 			auth					// 认证信息
}


type Response struct {
	Code	int `json:"code"`
	Msg		string	`json:"msg"`
	Data	messagePackage `json:"packages"`
}

func (cxt *Context) Response (resp *Response) {
	res,_ := json.Marshal(resp)
	cxt.ConnSocket.Write(res)
}