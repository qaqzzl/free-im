package service

import (
	"fmt"
	db "im/server/model"
)

// client auth handle
func (ctx *Context) ClientAuth() {
	//认证 ctx.Message.AccessToken	&& ctx.Message.UserID

	//记录链接信息
	fmt.Println(len(db.Connects[ctx.Message.UserID]))
	for v := range db.Connects[ctx.Message.UserID]  {
		fmt.Println(v)
	}
	//db.Connects[ctx.Message.UserID] = &db.Connect{
	//	Conn: &ctx.Conn,
	//	IsSignIn:true,
	//}
	ctx.Conn.Write([]byte("1"))
}

//client send message handle
func (ctx *Context) ClientSendMessage() {
	fmt.Println("client send message")
}

//client pull message handle
func (ctx *Context) ClientPullMessage() {

}