package socket

import (
	"encoding/json"
	"net"
)

type Context struct {
	ConnSocket		net.Conn				// 底层socket连接
	isClosed 		bool
	closeChan 		chan byte  				// 关闭通知

	InChan 			chan *[]byte			// 读队列 (入)
	OutChan 		chan *[]byte 			// 写队列 (出)
}


type Response struct {
	Status	int `json:"status"`
	Code	int `json:"code"`
	Msg		string	`json:"msg"`
	Data	interface{} `json:"data"`
}

func (cxt *Context) Response (resp Response) {
	res,_ := json.Marshal(resp)
	cxt.ConnSocket.Write(res)
}