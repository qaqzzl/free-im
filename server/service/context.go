package service

import (
	"encoding/json"
	db "im/server/model"
	"net"
)

type Context struct {
	Conn	net.Conn
	Message db.Message
}


type Response struct {
	Code	int `json:"code"`
	Msg		string	`json:"msg"`
	Data	interface{} `json:"data"`
}

func (cxt *Context) Response (resp *Response) {
	res,_ := json.Marshal(resp)
	cxt.Conn.Write(res)
}